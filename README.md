# Earth Backend


## Intro

When I was writing RTMP protocol from scratch, it came to my mind that I could build social media.


Since I watch anime sometimes, it comes to my mind why I don't build social media, which is has a feature that Otaku's watching anime online together(streaming real-time) and could chat with each other during streaming and also has a community for each anime series that can chat with each other.

After some feasibility study, I find out that I needed to build social media with distributed DB(which I was wrong about it). I thought it would be exciting to learn and develop my social media with Cassandra, so I started learning Cassandra and how to make DB with Cassandra and the challenge that would come in this social media with Cassandra etc. 

When I was writing authentication with go gin and Cassandra DB, I realized that creating social media with Cassandra is not that easy. The first problem I was faced was we have few ORM for Cassandra, which is written in golang, and finally, I find out gocql and gocqlx(which is, I think, a wrapper around gocql), and the second problem was these ORM's are very basic which is hard to develop my social media with these ORMs. For example, these ORMs do not support auto migration when fields in struct modify, added, or deleted. So what I did was write a simple auto migration logic(which worked, but it has a few bugs). After that, the next problem I was countering was clustering. The first attempt of my development was with one Cassandra instance. Still, when I added two other nodes, I ran out of resources, especially RAM. Sometimes due to this problem, the kernel immediately invokes OOM killer, which shut down my Cassandra instance. Then, I started to limit resources for the Cassandra instance from a docker-compose file, but after some attempt, it was fixed. It ran, but nodes have an unusual behavior which is made development miserable, so due to these problems and limited resources, I decided to stop developing this social media(but in the future, I will definitely build social media like what I described), so **this repository is abandon**.


What lessons i learned from this project(albeit it abandon) are:

1. ***Do your feasibility study more efficiently and practically by writing some small code, not by reading just some document.***
2. How does Cassandra work and how to design tables and queries in Cassandra, and what challenges does Cassandra DB introduce to the project.
3. How to use reflection in golang, For example, to build auto migration or other management tools to manage your codes easily without executing commands manually every time.
4. How distributed DB and distributed systems work, After this project, I attached to distributed systems and distributed DB.


All of the above lessons made me happy and extended my knowledge, which became helpful in my carries.