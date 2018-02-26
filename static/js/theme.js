$(document).ready(function(){
    $('input.timepicker').timepicker({});
});

function playTone() {
    var audio = new Audio('static/audio/tone.mp3');
    audio.play();
    console.log("playingtone")
};
