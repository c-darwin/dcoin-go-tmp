		var lang = 1;
		var room = 0;
		var chatMessage;
		var decryptChatMessage;
		var chatEncrypted = 0;

		$('#sendToChat').bind('click', function () {
			chatMessage = $('#myChatMessage').val();
			if (chatEncrypted == 1) {
				$.post( 'ajax?controllerName=encryptChatMessage', {
					'receiver': $('#chatUserIdReceiver').val(),
					'message': $('#myChatMessage').val()
				}, function (data) {
					chatMessage = data.success;
					decryptChatMessage = $('#myChatMessage').val();
					sendToTheChat()
				}, 'JSON');
			} else {
				sendToTheChat()
			}

		});

		if (!Date.now) {
			Date.now = function() { return new Date().getTime(); }
		}
		function sendToTheChat() {

			var chatMessageReceiver =  $('#chatUserIdReceiver').val();
			var chatMessageSender =  userId;
			var status = 0;
			if (chatEncrypted == 1) {
				status = 1
			}
			var signTime = Math.floor(Date.now() / 1000);
			var e_n_sign = get_e_n_sign( $("#key").text(), $("#password").text(), lang+","+room+","+chatMessageReceiver+","+chatMessageSender+","+status+","+chatMessage+","+signTime, 'chat_alert');
			$.post( 'ajax?controllerName=sendToTheChat', {
				'receiver': chatMessageReceiver,
				'sender': userId,
				'lang': lang,
				'room': room,
				'message': chatMessage,
				'decrypt_message': decryptChatMessage,
				'status': status,
				'sign_time': signTime,
				'signature': e_n_sign['hSig']
			}, function (data) {

			});
		}

		function scrollToBottom() {
			var objDiv = document.getElementById("chatwindow");
			console.log(objDiv.scrollHeight-67-objDiv.scrollTop)
			if (objDiv.scrollTop == 0 || objDiv.scrollHeight-67-objDiv.scrollTop == objDiv.clientHeight) {
				objDiv.scrollTop = objDiv.scrollHeight;
			}
		}
		$(document).ready(function() {
			$.post( 'ajax?controllerName=getChatMessages&first=1&room='+room+'&lang='+lang, {}, function (data) {

				if(typeof data.messages != "undefined" && data.messages !="") {
					$('#chatMessages').append(data.messages);
					scrollToBottom();
				}

			}, 'JSON');
			var intervalID = setInterval( function() {
				$.post( 'ajax?controllerName=getChatMessages&room='+room+'&lang='+lang, {}, function (data) {
					//if(typeof data.messages != "undefined" && data.messages !="") {
						console.log("data.messages", data.messages);
						$('#chatMessages').append(data.messages);
						scrollToBottom();
					//}
				}, 'JSON');

				var objDiv = document.getElementById("chatwindow");
				console.log(objDiv.scrollHeight, objDiv.scrollTop, objDiv.clientHeight)

			} , 1000);
			intervalIdArray.push(intervalID);
		});

		function setReceiver(nick, receiverId){
			$('#myChatMessage').val(nick+", ");
			$('#chatUserIdReceiver').val(receiverId);
			$("#selectReceiver").css("display", "none");
			$("#myChatMessage").css("display", "inline-block");
			console.log("receiverId", receiverId)
		}


		$('#chatLock').bind('click', function () {
			if ($(this).attr('class') == "fa fa-lock") {
				$(this).attr("class", "fa fa-unlock");
				$("#myChatMessage").css("display", "inline-block");
				$("#selectReceiver").css("display", "none");
				$("#myChatMessage").css("background-color", "#fff");
				$("#myChatMessage").css("color", "#000");
				chatEncrypted = 0
			} else {
				$(this).attr("class", "fa fa-lock");
				if ($("#chatUserIdReceiver").val() == "0") {
					$("#myChatMessage").css("display", "none");
					$("#selectReceiver").css("display", "inline-block");
				}
				$("#myChatMessage").css("background-color", "#BC5247");
				$("#myChatMessage").css("color", "#fff");
				chatEncrypted = 1
			}
		});

