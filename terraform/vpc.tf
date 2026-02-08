data "aws_availability_zones" "available" {}

resource "aws_vpc" "main" {
    cidr_block = "10.0.0.0/16"
    enable_dns_support = true
    enable_dns_hostnames = true

    tags = {
        Name = "diplomacy-vpc"
    }
}

resource "aws_subnet" "public" {
    count = 2
    vpc_id = aws_vpc.main.id
    cidr_block = cidrsubnet(aws_vpc.main.cidr_block, 8, count.index) # /24s
    availability_zone = data.aws_availability_zones.available.names[count.index]
    map_public_ip_on_launch = true

    tags = {
        Name = "diplymacy-subnet-public-${count.index}"
    }
}

resource "aws_subnet" "private" {
    count = 2
    vpc_id = aws_vpc.main.id
    cidr_block = cidrsubnet(aws_vpc.main.cidr_block, 8, count.index + 10) # /24s
    availability_zone = data.aws_availability_zones.available.names[count.index]

    tags = {
        Name = "diplomacy-subnet-private-${count.index}"
    }
}

resource "aws_internet_gateway" "igw" {
    vpc_id = aws_vpc.main.id
    tags = {
        Name = "diplomacy-igw"
    }
}

resource "aws_route_table" "public" {
    vpc_id = aws_vpc.main.id
    tags = {
        Name = "diplomacy-public-rt"
    }
}

resource "aws_route" "public_internet" {
    route_table_id = aws_route_table.public.id
    destination_cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
}

resource "aws_route_table_association" "public" {
    count = 2
    subnet_id = aws_subnet.public[count.index].id
    route_table_id = aws_route_table.public.id
}

resource "aws_eip" "nat" {
    domain = "vpc"
    tags = {
        Name = "diplomacy-nat-eip"
    }
}

resource "aws_nat_gateway" "nat" {
    allocation_id = aws_eip.nat.id
    subnet_id = aws_subnet.public[0].id
    tags = {
        Name = "diplomacy-nat"
    }
}

resource "aws_route_table" "private" {
    vpc_id = aws_vpc.main.id
    tags = {
        Name = "diplomacy-private-rt"
    }
}

resource "aws_route" "private_outbound" {
    route_table_id = aws_route_table.private.id
    destination_cidr_block = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat.id
}

resource "aws_route_table_association" "private" {
    count = 2
    subnet_id = aws_subnet.private[count.index].id
    route_table_id = aws_route_table.private.id
}
