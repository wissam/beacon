# beacon
Beacon started as a small project to manipulate a lifx bulb to be a
notification for events and errors that occur on cloud infrastructure and cicd
pipelines. 
Now I am having ideas to evolve it even more maybe a multi platform
notification system? will see.


## Services


### Vislog

The visual logger, wrapper for lifx api to have functions for Error,Warning,
Success, etc.
This should support other brands in the future.


### Redis pubsub

The message queue mechanism to control vislog



### Hook Server 

Webhooks url generator that handle different type of services



### Frontend
I will not be working on this, I will let my brother have a go at it.

### Backend
User, Groups, Bulbs Management with a db.


...fun :)
