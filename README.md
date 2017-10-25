# MESS for Makers
### Member and Event Supervisory System for Makers
At MakeICT, we have gotten tired of fighting with our current system for managing members and events.  Lack of features, lack of support, broken API's, the list goes on.  So we've decided to make a better system, since we're makers and all.  This will give us the access we need to our users and events, and give us opportunities to implement new features like equipment management and to work with new developers.  We are in the very first stages of this project, so if you would like to be involved but don't have a lot of experience, now is a great time. Wee need developers at all levels, front end, back end, user experience, database, API and more.

The first steps are to get something minimal that can begin to take the place of the current system used by MakeICT.  In brief, that means it will need to be a web-based system that allows us to create events, reserve space, and allow users to register for those events.  We also would like to give users some more control over the profiles and data registered in the system, and give administrators more tools for reporting usage and metrics.  We have also decided to continue supporting those who wish to pay for classes and membership dues via cash or check, so we will need at least a minimal invoicing system at the start.  For a more complete list of features, including future wishes, check out the [Google Doc](https://docs.google.com/document/d/1kCKM_0OuQ-ox3oTD7ylt77YPgt1ZrhlLrgR1eQ0qVwc/edit).  

We also hope to make this project about a lot more than just MakeICT.  We feel that every makerspace could use a good management tool for free. We'd like to keep in mind a design philosophy that makes it easy for other organizations to deploy, customize, and extend this platform.

# Resources 
You can join the discussion on Slack, or come to a scheduled meeting.  Check out [devICT's home page](http://devict.org) for information on joining Slack and for joining the Meetup for event notifications.

Our current design mockups are on [NinjaMock](https://ninjamock.com/s/JC7Q9).  Check them out to see the direction we are taking, especially if you would like to help with the front-end development.

The list of requirements for the software is in this [Google Doc](https://docs.google.com/document/d/1kCKM_0OuQ-ox3oTD7ylt77YPgt1ZrhlLrgR1eQ0qVwc/edit).  

//TODO
Expand this with getting started resources for the different technologies used in this project.

# Contributing
Contributions are welcome, you can start by reading this section for some of the best practices and methods for managing this project.
Next, check the projects or issues for something that interests you, or that you have skills that you could offer.
When you find something, see if there is a branch for that feature or bug, if not make a new one.  

### What technology does this prjoect use?
We are currently building the server backend in Go (1.7.3), the frontend with [Material.io](https://material.io), and using PostgreSQL (9.3) for storage.  Various other libraries may be used, but those will be explained in the comments in the files that use them.

### How should I add or work on a feature?
Check the projects or issues to find something you want to work on.  Then, check the branches to see if someone is already working on it.  If they are, you can join that branch, or if not, you will need to create a new branch. More information on the git Feature Branch workflow can be found in [this tutorial](https://www.atlassian.com/git/tutorials/comparing-workflows#feature-branch-workflow) by Atlassian.  It also has a very good beginners guide to using Git if you never have before.
After you have worked on your feature and have it working you can submit a pull request and after your code is reviewed it will be merged into the main branch.

### How will all this code be tested?
At this time, with very little written in the way of code or tests, the plan is to use acceptance testing to make sure that routes, validation, server responses, etc are working.  Unit tests may be used in some cases to prove that specific bugs are fixed and prevent regression.

# Setting up with Cloud9
Cloud9 has PostgreSQL, Git, and Go already set up and is a workable Go IDE, so it's easy for someone new to Go or software development to use for getting started. This guide does assume you are familiar with git, and creating and navigating directories on Linux.
1. Go to [Cloud9](https://c9.io)
2. Login with your github account
3. Create new blank workspace
4. Create 3 new folders, `bin`, `pkg`, and `src` for your compiled binaries, compiled libraries, and source code respectively.  These are the folders that Go expects by default.
5. `$GOPATH` is set up automatically, but you can type `c9 open ~/.profile` and add this at the end of the file to make running binaries easier. 
```
export PATH=$PATH:$GOPATH/bin
```
6. In the same `.profile` add the following to automatically run the PostgreSQL server: You will have to close and reopen the terminal for these to take effect.
```
function checkstart {
	service=$1
	if [[ ! $(ps -ef | grep -v grep | grep "$service" | wc -l) > 0 ]]
	then
		sudo service $service start &
	fi
}
checkstart postgresql
```
7. Create a postgres user by typing `psql` and then 
```
CREATE ROLE <username> WITH LOGIN PASSWORD ‘<password>’ CREATEDB;
```
8. Create a postgres database while still in psql with 
```
CREATE DATABASE <database> OWNER <username>;
```
9. Create folder `src/github.com/makeict/MESSforMakers` and change to that folder.
From the folder, run `git clone https://github.com/MakeICT/MESSforMakers.git`
10. You then need to install all the build dependencies with 
```
go get github.com/jmoiron/sqlx github.com/gorilla/sessions github.com/gorilla/mux github.com/justinas/alice github.com/lib/pq
```
  - This list is subject to probably a lot of change. If you get errors that a library cannot be found, just `go get` that library
11. Prepopulate the database with 
```
psql postgres://<username>:<password>@localhost:5432/<database> -f test_tables.sql
```
12. Configure the server by opening `config.json` and setting the username, password, and database to whatever you chose earlier, and the host and port to `localhost` and `5432`.
13. At this point you should be able to type `go install` and then `MESSforMakers` and the server should run. You can then click the “Preview” button at the top of the editor to see the application running in a browser window.
14. All these instructions should work on any linux distro, but if you don’t use Cloud9 you will have to set up Git, PostgreSQL, and Go manually.

