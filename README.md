# Instagram-Backend-API
A Backend API for a normal Social Media Platform

## External Dependencies

- mongo-driver 
- godotenv (to handle secrets using environment variables, not required)

## List Of Endpoints

#### ``` /users ```
creates a user and encrypts the password using ciphers before storing

#### ``` /users/?id=<id here> ``` 
fetches user details of given id

#### ``` /posts ``` 
creates a post using data from the POST request's body

#### ``` /posts/<id> ``` 
fetches post details for the given id

#### ``` /posts/users/?id=<id here>&limit=<pagination limit>``` 
fetches posts of the user with given id within a limit 



## Pagination

This API preloads a set number of filtered posts beforehand to increase performance and incorporate Lazy Loading Concept. It also returns a last_id for the last post loaded so that the nest set number of posts can be loaded.

![](https://i.imgur.com/9Xsa8bp.png)



## Screenshots

- Creating a User with password encryption

![](https://i.imgur.com/hSEObRV.png)

- User created in mongo
![](https://i.imgur.com/IV3rjwT.png)

- Fetching a User

![](https://i.imgur.com/S6mVjPR.png)

- Create a Post

![](https://i.imgur.com/u6OMQIi.png)

- Post created in mongo

![](https://i.imgur.com/7yQTgwG.png)

- Fetch a single post

![](https://i.imgur.com/IMnLK8n.png)

- Fetch posts created by a user with pagination

![](https://i.imgur.com/9Xsa8bp.png)
