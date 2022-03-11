# ADR-010-database

* Creation Date: 10/03/2022
* Status: Draft | Pending | Accepted | Denied

## Context
The Voogle project is a demonstrator for Squarescale. We need to store video informations in a persistent system (video sources are stored in a Bucket S3). Squarescale provide support for 3 database engine : MariaDB, MySQL and PostgreSQL. This database aims to be simple and short.

## Decision
We will use the option 1 with MariaDB.

## Options

### 1. MariaDB
*"MariaDB is a community-developed, commercially supported fork of the MySQL relational database management system (RDBMS), intended to remain free and open-source software under the GNU General Public License."*
[MariaDB]https://en.wikipedia.org/wiki/MariaDB
It's a MySQL Fork intended to remain free and open-source. It then maintains high compatibility with MySQL up to 7th version of MySQL include.

#### Benefits
- Used more and more
- Most Linux distributions use MariaDB instead of MySQL as default DBMS.
- Light : good for small database and small machines
- Community-developed (high number of new features)
- Support for MariaDB with Golang : https://mariadb.com/fr/resources/blog/using-go-with-mariadb/
- Allow partitioning (query performance and availability)

#### Drawbacks

### 2. MySQL
*"MySQL is a free and open-source RDBMS under the terms of the GNU General Public License, and is also available under a variety of proprietary licenses.*
*MySQL has stand-alone clients that allow users to interact directly with a MySQL database using SQL, but more often, MySQL is used with other programs to implement applications that need relational database capability."*
[MySQL]https://en.wikipedia.org/wiki/MySQL

#### Benefits
- Massively used
- Light : good for small database and small machines
- Allow partitioning (query performance and availability)
  
#### Drawbacks
- Free version limited

### 3. PostgreSQL
*"PostgreSQL is a free and open-source RDBMS emphasizing extensibility and SQL compliance.*
*PostgreSQL features transactions with ACID properties, automatically updatable views, materialized views, triggers, foreign keys, and stored procedures.It is designed to handle a range of workloads, from single machines to data warehouses or Web services with many concurrent users."*
[PostgreSQL]https://en.wikipedia.org/wiki/PostgreSQL

#### Benefits
- Massively used
- Faster than others for substantial data
- Lot of advanced features

#### Drawbacks
- Heavy for small database
- Only supports Master-Slave replication
- Table partitioning not supported

## Technical resources
-  mariadb-vs-postgresql https://www.ionos.fr/digitalguide/hebergement/aspects-techniques/mariadb-vs-mysql/
-  mariadb-vs-postgresql https://hevodata.com/learn/mariadb-vs-postgresql/
