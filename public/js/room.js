// ---- P2P ------
var roomID = $('#myroom').val();
var connectedIds = [];
var connectedPeers = {};
var initiator = false;
var peerReady = false;

function initialize() {
  initiator = true;
  var data = JSON.stringify({"type": "initialize", "roomID": roomID, "uid": uid});
  socket.send(data);
  peerReady = true;
}

function leave() {
  for (var i=0; i < connectedIds.length; i++) {
    // send closing signal to others
    var sigTo = connectedIds[i];
    var data = JSON.stringify({"type": "close", "to": sigTo, "uid": uid});
    socket.send(data);
  }
  setTimeout("backToDashboard()", 300);
}

function backToDashboard(){
    location.href='/dashboard';
}

function createConnection(sigId, initFlg) {
  console.log("Connect to: " + sigId);

  var newPeer = new SimplePeer({ initiator: initFlg, stream: localStream });
  connectedPeers[sigId] = newPeer;
  sendSignal(newPeer, sigId);
}

function sendSignal(peer, id) {
  peer.on('signal', function (data) {
    data.to = id;
    data.uid = uid;
    data.roomID = roomID;
    var text = JSON.stringify(data);
    socket.send(text);
    console.log("Signal send to: " + id + text);
  })
  peerReady = true;
}

function receive(signal, peer, sigFrom) {
  peer.signal(signal);
}

function startCall(peer, id) {
  if ($('#start_video_call').length == 1) {
    // change disploayed button
    $('#start_video_call').remove();
    $('#start_voice_call').remove();
    leave_button = '<div class="col-xs-6" id="leave_room"><div class="form-group call-form"><button class="btn btn-block btn-lg btn-danger" type="button" onclick="leave();">Leave Room</button></div></div>';
    $('#room_title').after(leave_button);
  }

  // start dictation
  startDictation();
  // start video call
  $("#videos").append("<video id='video-"+id+"' autoplay style='width: 240px; height: 180px; border: 1px solid black;'></video>")
  peer.on('stream', function (remoteStream) {
    $("#video-"+id).prop('src', window.URL.createObjectURL(remoteStream));
  })
  peer.on('error', function (err) {
    console.log("Peer Error has happened");
    closeCall(peer, id)
  })
  peer.on('close', function () {
    console.log("Peer has closed " + peer.channelName);
    closeCall(peer, id)
  })
}

function closeCall(peer, id) {
  connectedIds = connectedIds.filter(function(v){
    return v != id;
  })
  if (connectedIds.length == 0) {
    // stop dictatoin
    stopDictation()
  }

  // close video call
  $("#video-"+id).remove();
  peer.destroy()
}

// ---- socket ------
var domain = $('#mydomain').val();
var uid = $('#myid').val();
var uname = $('#myname').val();
var uri = 'wss://' + domain + ':443/ws';
var socket = null;

function socketOpen() {
  if (socket == null) {
    // WebSocket の初期化
    socket = new WebSocket(uri);
    socket.onopen = onOpen;
    socket.onmessage = onMessage;
    socket.onclose = onClose;
    socket.onerror = onError;
  }
}
socketOpen();

function onOpen(event) {
  console.log('socket opened.');
  registerSocket();
}

function onMessage(event) {
  var signal = JSON.parse(event.data);
  console.log("Received.");
  console.log(signal);

  if (signal.type == 'info') {
    // confirm to register socket client
    console.log("Socket is registerd. MyID is :" + signal.uid);
  } else if (signal.type == 'initialize') {
    var pool = signal.ids;
    for (var i=0; i < pool.length; i++) {
      var sigTo = pool[i];
      // まだ繋がっていなければ接続
      if (connectedIds.indexOf(sigTo) == -1) {
        createConnection(sigTo, true)
      }
    }
  } else if (signal.type == 'conversation') {
    // 会話のログ記述
    var u = signal.uname;
    var c = signal.content;
    $("#conversation").append("<p>"+u+": "+c+"</p>");
  } else if (signal.type == 'offer') {
    // だれから送られてきたのか取得
    var sigFrom = signal.from;
    delete signal['from'];
    delete signal['room'];
    createConnection(sigFrom, false)
    receive(signal, connectedPeers[sigFrom], sigFrom);
    connectedIds.push(sigFrom);
    startCall(connectedPeers[sigFrom], sigFrom);
  } else if (signal.type == 'answer') {
    // だれから送られてきたのか取得
    var sigFrom = signal.from;
    delete signal['from'];
    delete signal['room'];
    receive(signal, connectedPeers[sigFrom], sigFrom);
    connectedIds.push(sigFrom);
    startCall(connectedPeers[sigFrom], sigFrom);
  } else if (signal.type == 'close') {
    // だれから送られてきたのか取得
    var sigFrom = signal.from;
    closeCall(connectedPeers[sigFrom], sigFrom);
  } else {
    // type: candidate
    // だれから送られてきたのか取得
    var sigFrom = signal.from;
    delete signal['from'];
    delete signal['room'];
    receive(signal, connectedPeers[sigFrom], sigFrom);
  }
}

function onError(event) {
  //console.log("Socket Error");
}

function onClose(event) {
  console.log("Socket is Closed.");
  socket = null;
  socketReady = false;
  // reopen after 3 seconds
  setTimeout("socketOpen()", 3000);
}

function registerSocket() {
  var data = JSON.stringify({ "uid": uid, "type": "register"});
  socket.send(data);
}

// ---- Video ----
var localStream = null;

window.onload = function(){
  navigator.getUserMedia  = navigator.getUserMedia ||
                            navigator.webkitGetUserMedia ||
                            navigator.mozGetUserMedia;
  window.URL = window.URL || window.webkitURL;

  var video = document.getElementById('my-video');
  navigator.getUserMedia({video: true, audio: true}, function(stream) { // for success case
    localStream = stream;
    video.src = window.URL.createObjectURL(stream);
   },
   function(err) { // for error case
    console.log(err);
   }
  );
}

// 音声認識
var WToken = $('#watsonToken').val();
var stream = WatsonSpeech.SpeechToText.recognizeMicrophone({
    token: WToken,
    model: 'ja-JP_BroadbandModel'
});

function startDictation() {
  console.log('start Dictation!!!')

  stream.setEncoding('utf8'); // get text instead of Buffers for on data events

  stream.on('data', function(data) {
    console.log(data);
    $("#conversation").append("<p>"+uname+": "+data+"</p>");

    // send others
    for (var i=0; i < connectedIds.length; i++) {
      var data = JSON.stringify({"type": "conversation", "content": data, "uname": uname, "to": connectedIds[i]});
      socket.send(data);
      //console.log(socket);
      //console.log("send conversation");
    }
  });

  stream.on('error', function(err) {
    stream.stop.bind(stream);
  });
}

function stopDictation() {
  stream.stop.bind(stream);
}
