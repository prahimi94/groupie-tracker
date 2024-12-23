# groupie-tracker

## Description
Groupie Trackers consists on receiving a given API and manipulate the data contained in it, in order to create a site, displaying the information.

This project used to filter artists data based on user selections.

It will be given an API, that consists in four parts:

The first one, artists, containing information about some bands and artists like their name(s), image, in which year they began their activity, the date of their first album and the members.

The second one, locations, consists in their last and/or upcoming concert locations.

The third one, dates, consists in their last and/or upcoming concert dates.

And the last one, relation, does the link between all the other parts, artists, dates and locations.

Given all this you should build a user friendly website where you can display the bands info through several data visualizations (examples : blocks, cards, tables, list, pages, graphics, etc). It is up to you to decide how you will display it.

This project also focuses on the creation of events/actions and on their visualization.

The event/action we want you to do is known as a client call to the server (client-server). We can say it is a feature of your choice that needs to trigger an action. This action must communicate with the server in order to recieve information, ([request-response])(https://en.wikipedia.org/wiki/Request%E2%80%93response)
An event consists in a system that responds to some kind of action triggered by the client, time, or any other factor.


[Base Api Address:](https://groupietrackers.herokuapp.com/api)

---

## Authors
- **[Parisa Rahimi Darabad]** - Senior Backend Developer  
- **[Majid Rouhani]** - Senior Backend Developer  

---

## Usage
### How to Run
1. Clone this repository to your local machine:
   ```bash
   git clone https://01.gritlab.ax/git/prahimi/groupie-tracker.git

2. Navigate to the project directory:
    ```bash
    cd groupie-tracker

3. Start the development server:
    ```bash
    go run .

4. Open your browser and navigate to:
    ```arduino
    http://localhost:8080

5. Go to artists menu and filter your selections.

6. You can run test from root with this command:
    ```arduino
    go test ./...

    
## Project Structure and Implementation
Project has 2 main components

Backend: Include Dockerfile, webserver, Api and Tests

Frontend: Include html templates, error files and assets

In the root of the project there is a docker-compose.yml file to run project from the root with docker while keep structure of the project organized.