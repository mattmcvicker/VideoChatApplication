# VideoChatApplication

## Project link:
[Proposal](https://docs.google.com/document/d/1U90b1wiZy8jIflSR2nrX0vsznA1GCBe5P63_CAlR3Zc/edit?pli=1)

### Project Description
Our target audience is people who want to have an open debate about specific subjects with people that they disagree with. These subjects can be about anything, from politics to sports to food. A secondary audience that we could see using this application is TA’s using this application as an icebreaker for students. We decided to implement a video chat instead of something like a text chat because it's a lot easier to communicate and get your ideas across when you’re talking face-to-face with someone compared to when you’re trying to communicate with them through text. This is also a way to weed out trolls, making it so that there is some amount of accountability for the things that you’re saying rather than it being behind an anonymous text chat. 
Our audience wants to use this application because they want to learn more about other sides of issues that they care about. It’s very easy to get into an echo chamber when it comes to important topics, so by actually talking with someone who has an opinion different from their own, our audience can educate themselves on the opposite side, and even learn more about their own stance. 
	We want to build this application because we feel that it could be a good way for people who have differing ideologies on a topic to understand the underlying thoughts that a person with a different opinion on the subject matter might have. It would also use a lot of the techniques we have learned in this class.

### Technical Description (Diagram in link above)
| Priority | User                                | Description |
|----------|-------------------------------------|-------------|
| P0       | As a user                           | I want to be able to sign into my account with all my topics and opinions already saved|
| P0       | As a user                           | I want to be able to change the opinions I have set on my profile|
| P0       | As a user                           | I want to be able to find other people with differing opinions on a specific topic|
| P0       | As a user                           | I want to be able to look into new topics that I have not yet set and set my opinions on those topics|
| P2       | As a user who is neutral on a topic | I want to be able to get matched with people who feel strongly about a particular side to educate myself on that topic|
| P2       | As a teacher/user                   | I want to be able to create my own topics and prompts|


Users will be able to sign into their account, which is stored in the MySQL database. Their opinion on topics is also stored in the MySQL database.
Users will be able to change the opinions they have set on their profile, which are stored in the MySQL database.
Users will be assigned a video room id, based on the topic they want to debate which is handled by WebRTC and our Firestore. WebRTC will create a video room id and send it to the firestore, and we will handle matching users and pairing them with a connection.
Users will be able to search topics, or view some of the more popular topics on the website and set their opinions on those, which will be stored in the MySQL database.
Users will be able to set their opinion on a particular topic as “neutral,” and that will be saved to their profile in the MySQL database. Then people will get matched to people who feel strongly about one side or other of the topic.
Users will be able to create their own topics, asking different questions for different opinions, and have people answer those questions and get matched accordingly. These topics will be stored in the MySQL database.

### Endpoints: 
See Link Above

### MYSQL:
See Link Above







