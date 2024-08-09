# Defining the project

## What is the project?
This project (Baasa) is what is a 'Backend as a service', it basicly a server that has a database and have an authentication system.
Clients can send requests to the server like accessing the server data and logging into their user or creating a new one.
The admin can configure which data the clients can access and how.

## What is the MVP (Minimal Viable Product)?
The admin having a dashboard when he can create and edit the collections of the server.
The admin can create and edit new queris - which a basicly a custom query from the database, clients can send a request to each query and can the query resultsback from the database.
Each query have access rules which is basicly the rules of who can call the query, the admin can edit the access rules from the dashboard.
The clients can access each collection in 5 different methods and each method has its own custom access rules.
The methods are:
    - Search - searching in the collection with custom filters.
    - View - getting a singe record from the collection by searching with custom filters.
    - Create - creating a new record in the collection.
    - Update - updating a record in the collection by searching with custom filters.
    - Delete - updating a record in the collection by searching with custom filters.

In the database there would be a 'users' collection which would hold all the information of the users.
The client can register and login to server.
When logging into the server the client would get a session in return which he can use in he next requests to use his user without logging in again.
There should be an sdk in javascript to access the server and all its features.

## What are the nice to haves?
Adding file storage to the baas.
Having a nice UI in the admin dashboard.
Having email validation when registering.
Having SDK in other programing languages aside from js.
...

## When will the project will be complete?
When the MVP will be complete.

---
