'use strict'
let localStorage = window.localStorage

//init button event listeners
let topicmain = document.querySelector("#topics-body")
let quizmain = document.querySelector("#quiz-body")
let formsection = document.querySelector("#topic-form-section")
let listsection = document.querySelector("#topic-list")
let listsection2 = document.querySelector("#topic-list-section")
let currentTopicID = -1
let currentTopicName = ""


//Post a new topic Event listeners
document.querySelector("#postbtn").addEventListener("click",
    () => {
        toggleVisibilty()
    }
);

document.querySelector("#submit-button").addEventListener("click",
    () => {
        toggleVisibilty()


        //post the new topic
        //generatetopiclist()
    }
);
///////////


document.querySelector("#form-back-button").addEventListener("click",
    () => {
        toggleVisibilty()
    }
);


//button helpers
function toggleVisibilty() {
    listsection.classList.toggle("d-none")
    listsection2.classList.toggle("d-none")
    formsection.classList.toggle("d-none")
    document.querySelector("#postbtn").classList.toggle("d-none")
}

function generateTopicList(){
    listsection.innerHTML = "";
    //placeholder

    //set date to right format (got code from somewhere on the internet)
    var today = new Date();
    var dd = String(today.getDate()).padStart(2, '0');
    var mm = String(today.getMonth() + 1).padStart(2, '0'); //January is 0!
    var yyyy = today.getFullYear();

    today = mm + '/' + dd + '/' + yyyy;

    
    for (var i = 0; i < 5; i++) {
        listsection.appendChild(topicitem(1, "Salads? "+i , today, "squidward69", 69))
    }

    //fetch topics, then for each topic,
    // response.forEach(
    //     (item) => {
    //         listsection.appendChild(topicitem(item.id, item.name, item.createdAt, item.creator, item.votes))
    //     }
    // )
}



//on page init
generateTopicList()


///

//create a topic item DOM object
function topicitem(id, topicname, timecreated, topicauthor, likes) {
    let frame = document.createElement("div")
    frame.classList.add("container", "jumbotron", "bg-dark", "text-light");
    
    //name
    let titleframe = divWithClass(["row"])
    let titlebody = divWithClass(["col"])
    let titlechild = document.createElement("h3")
    titlechild.textContent = topicname

    titlebody.appendChild(titlechild)
    titleframe.appendChild(titlebody)

    //info
    let infoframe = divWithClass(["row"])
    let info = [divWithClass(["col"]), divWithClass(["col"]), divWithClass(["col"])]
    info[0].textContent = "Created by: " + topicauthor
    info[1].textContent = "at: " + timecreated.toString()
    info[2].textContent = "Likes: " + likes

    info.forEach((item)=>{infoframe.appendChild(item)})

    //buttons
    let buttonframe = divWithClass(["row"])
    let buttons = [
        buttonGen("Enter", "col-7 btn m-1 btn-primary", id),
        buttonGen("Like", "col-2 btn m-1 btn-light", id),
        buttonGen("Dislike", "col-2 btn m-1 btn-light", id)
    ];

    //ENTER btn event listener
    buttons[0].addEventListener("click", 
        () => {
            document.querySelector("#topicname").textContent = topicname
            topicmain.classList.add("d-none")
            quizmain.classList.remove("d-none")
            quizBtnEvents(id)
            //do stuff with quizzes
        }
    );
    
    //LIKE btn event listener
    buttons[1].addEventListener("click", 
        () => {
            //on click:
            //PATCH topic
            //generateTopicList()
        }
    );
    
    //DISLIKE btn event listener
    buttons[2].addEventListener("click", 
        () => {
            //on click:
            //PATCH topic
            //generateTopicList()
        }
    );

    buttons.forEach((item)=>{buttonframe.appendChild(item)})
    
    frame.appendChild(titleframe)
    frame.appendChild(infoframe)
    frame.appendChild(buttonframe)

    return frame
}


//helper methods
function getFormattedDate(date) {
    var year = date.getFullYear();
  
    var month = (1 + date.getMonth()).toString();
    month = month.length > 1 ? month : '0' + month;
  
    var day = date.getDate().toString();
    day = day.length > 1 ? day : '0' + day;
    
    return month + '/' + day + '/' + year;
  }

//for both these methods, c is the desired html class
function divWithClass(c) {
    let div = document.createElement("div")

    c.forEach(
        (item) => {
            div.classList.add(item)
        }
    )

    return div
}

function buttonGen(content, c) {
    let btn = document.createElement("button")

    //i forgot that this method existed lol
    btn.setAttribute("class", c)

    btn.textContent = content
    return btn
}

// Quiz


function quizBtnEvents(topicID) {
    let truebtn = document.querySelector("#truebtn")
    let falsebtn = document.querySelector("#falsebtn")

    truebtn.addEventListener("click",
        () => {
            
            //add to queue as true for topic topicid
        }
    );

    falsebtn.addEventListener("click",
        () => {
            //add to queue as false for topic topicid
        }
    );

}