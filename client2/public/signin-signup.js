(function() {
    "use strict";

    //functionality for user signins/signups
    const apibase = "https://api.kenmasumoto.me";
    const userHandler = "/v1/users";
    const myUserHandler = "/v1/users/me";
    const sessionHandler = "/v1/sessions";
    var user = "";

    window.addEventListener("load", initialize);

    var signupPage = document.getElementById("signup");
    var signinPage = document.getElementById("signin");
    var queuePage = document.getElementById("queue");
    var videoPage = document.getElementById("video-grid");
    var topicsPage = document.getElementById("topics-body");

    function initialize() {
        // check if there is a user currently signed in
        getCurrentUser().then( () => {
            // define event listeners
            document.getElementById("signup-form").addEventListener("submit", signup);
            document.getElementById("signin-form").addEventListener("submit", signin);
            document.getElementById("signup-instead-btn").addEventListener("click", showSignup);
            document.getElementById("queue-button").addEventListener("click", function(){
                videoPage.classList.remove("hidden");
                queuePage.classList.add("hidden")
                startCall();
            });
        })
    }

    ////////// define functions /////////

    async function getCurrentUser() {
        const authToken = localStorage.getItem("Authorization");
        if (authToken == null || authToken == "") {
            return;
        }
        const response = await fetch(apibase + myUserHandler, {
            headers: new Headers({
                "Authorization": localStorage.getItem("Authorization")
            })
        });
        if (response.status > 300) {
            alert("Unable to verify login. Logging out.");
            localStorage.setItem("Authorization", "");
            showSignup()
            return;
        }else {
            user = await response.json();
            topicsInit();
            showTopics();
        }
    }

    async function signup(e) {
        e.preventDefault();
        const email = document.getElementById("signup-email").value;
        const username = document.getElementById("signup-username").value;
        const fname = document.getElementById("signup-fname").value;
        const lname = document.getElementById("signup-lname").value;
        const pass = document.getElementById("signup-password").value;
        const passConf = document.getElementById("signup-passwordConf").value;

        let postData = {
            "email": email,
            "userName": username,
            "firstName": fname,
            "lastName": lname,
            "password": pass,
            "passwordConf": passConf
        }

        const response = await fetch(apibase + userHandler, {
            method: "POST",
            body: JSON.stringify(postData),
            headers: new Headers({
                "Content-Type": "application/json"
            })
        });
        if (response.status >= 300) {
            const error = await response.text();
            console.log("error posting user: " + error);
            return;
        }

        // set auth token
        const authToken = response.headers.get("Authorization");
        localStorage.setItem("Authorization", authToken);

        // get user 
        user = await response.json();

        topicsInit();
        showTopics();
    }

    async function signin(e) {
        e.preventDefault();
        const email = document.getElementById("signin-email").value;
        const pass = document.getElementById("signin-password").value;

        let postData = {
            "email": email,
            "password": pass
        }

        const response = await fetch(apibase + sessionHandler, {
            method: "POST",
            body: JSON.stringify(postData),
            headers: new Headers({
                "Content-Type": "application/json"
            })
        });
        if (response.status >= 300) {
            const error = await response.text();
            console.log("error posting user: " + error);
            return;
        }

        // set auth token
        const authToken = response.headers.get("Authorization");
        localStorage.setItem("Authorization", authToken);

        // get user 
        user = await response.json();
        topicsInit();
        showTopics();

    }

    async function showSignin() {
        signupPage.classList.add("hidden");
        signinPage.classList.remove("hidden")
    }

    async function showSignup() {
        signinPage.classList.add("hidden");
        signupPage.classList.remove("hidden");
    }

    async function showTopics() {
        signinPage.classList.add("hidden");
        signupPage.classList.add("hidden");
        topicsPage.classList.remove("hidden");
    }

    async function showVideos() {
        queuePage.classList.add("hidden");
        videoPage.classList.remove("hidden");
        startCall()
    }

})();


