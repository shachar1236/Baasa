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
# Project components
- Admin dashboard
- Database
- Client API.
- Access Rules to the database.
- User authentication.
---
# Project features
## Admin dashboard
1. Show collections.
2. Creating new collection.
3. Changing collection properties.
4. Deleting collection.
5. Show collection records.
6. Edit collection record.
7. Creating new record in collection.
8. Edit collection access rules for each of the 5 access methods.
9. Show custom queris.
10. Creating new custom query.
11. Editing query.
12. Deleting query.
13. Editing query access rules.
## Database
1. Create user.
2. Get user/s.
3. Update user.
4. Delete user.
5. Create collection.
6. Get collection/s.
7. Update collection.
8. Delete collection.
9. Create query.
10. Get query.
11. Update query.
12. Delete query.
13. Query custom query.
14. Search in database with custom filters. (Same as Client API #_)
15. 
---
