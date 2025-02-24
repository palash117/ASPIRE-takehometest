Twitter apis 
## Architecture:
The system is divided into 4 modules which right now are written in all together in a same service. Idealy the modules should be deployed seperately so that we can scale them independenly of each other. 

1)  Http Controller: Responsible for serving all the http requests.
2)  User Module: Responsible for storing and maintaining User data and user relationship mappings. 
3)  Post Module: Responsible for storing posts created by users
4)  Timeline Module: Responsible for maintaining the timelines for all users. 

### Design decisions: 
1)  Monolith: all modules are currently in same repo, ideally should be separated as discussed above
2) Single Db: Currently I am using a single db ( SQLite) for all the above mentioned modules, in production, the data bases need to be separate  ,as the type of queries being supported are very different
	a) UserModule: for production we should use a sql based db for maintaining user table. For maintaining the user relationship a graph based db would be best.
	b) Post Module: For prod either of sql/nosql dbs should suffice
	c) Timeline Module: Since the data will be multiple orders of post , and the can be easily paritionable based on the user id, we can use a Nosql based distributable db for better availability .
3) Popular users: For popular users, (with some follower count > x) the post details is not duplicated for each user in their timelines, 
   Instead , when we are fetching the timeline, we query the posts table seperately to get posts from these popular users . This is done to maintain a balance between database read and write. 

## Table structure
- User : Maintains user related data
- Follower : Maintains user to user folloings
- Post : Contains post entries created for every user
- Timeline ( user_feed is typo in the image , please read as timeline):  
	Table containing timeline contents for users.  The data is duplicated from posts. 
	The rows are indexed by its primary key and also by user_id column. This is to speed up the queries for fetching all the rows for a user. 
	Since the application is supposed to be read heavy for normal usecases, When ever a Post is created by a user , a row for each user is created in timeline table containing the post details.
![[https://raw.githubusercontent.com/palash117/ASPIRE-takehometest/refs/heads/main/tableStructure.png]]
![[sequenceDiagram.png]]
### Api list:
- Post 
	- Post /twitter/{user_id}/post
		
- Get timeline for user
	- Get /twitter/{user_id}/timeline

#### _Helper apis_

- Create user   
	- Post /twitter/user

- Get all users 
	- Get /twitter/user/all

- Follow user
	- Post /twitter/{user_id}/follow
	
- Get followers for user 
	- Get /twitter/user_id/followers

- Get users following a user
	- Get /twitter/user_id/following


### Curls to generate data:

#### Create users:
```
curl --location 'http://localhost:8080/twitter/user' \
--header 'Content-Type: application/json' \
--data-raw '{
"userName":"Frieren",
"email":"Frieren@mail.com" 
}'

curl --location 'http://localhost:8080/twitter/user' \
--header 'Content-Type: application/json' \
--data-raw '{
"userName":"Fern",
"email":"fern@mail.com" 
}'

curl --location 'http://localhost:8080/twitter/user' \
--header 'Content-Type: application/json' \
--data-raw '{
"userName":"Himmel",
"email":"himmel@mail.com" 
}'

```

We have created 3 users , Frieren (user_id 1), Fern (2), Himmel(3)
#### Create User Folloings:
```
curl --location 'http://localhost:8080/twitter/user/1/follow' \
--header 'Content-Type: application/json' \
--data '{
	"followUserId": 2
}'

curl --location 'http://localhost:8080/twitter/user/1/follow' \
--header 'Content-Type: application/json' \
--data '{
	"followUserId": 3
}'

curl --location 'http://localhost:8080/twitter/user/2/follow' \
--header 'Content-Type: application/json' \
--data '{
	"followUserId": 3
}'

```

Now we have 3 relations
Frieren (1) ->(follows) -> Fern(2)
Frieren (2) ->Himmel(2)
Ferm -> Himmel

#### Create Post 
```
curl --location 'http://localhost:8080/twitter/3/post' \
--header 'Content-Type: application/json' \
--data '{
    "subject": "My favourite flower",
    "contents":"The blue-moon flower from my home town is my favourite"
}'

```

#### Check Timeline
```
curl --location 'http://localhost:8080/twitter/1/timeline?pageNo=0'

#######################
{"timeline":{"Posts":[{"Id":26,"AuthorId":3,"Subject":"My favourite flower","Contents":"The blue-moon flower from my home town is my favourite","CreatedAt":"2025-02-24T18:20:15.213123879Z"}]},"famousPosts":[],"pageNo":0}

```

the curl returns the feed for Frieren (1) which contains the post from  Himmel(3)

## Key notes
For users with follower count > 3 ( currently hardcoded due to lack of time)
the newly created posts from these users are not duplicated in the timeline table for all following users. 
