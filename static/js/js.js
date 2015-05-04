
function file_upload (file_id, progress, type, script) {
    var
        $f = $('#'+file_id),
        $p = $('#'+progress),
        up = new uploader($f.get(0), {
            url:'ajax/'+script,
            prefix:'file',
            type:type,
            progress:function(ev){ $p.html(((ev.loaded/ev.total)*100)+'%'); $p.css('width',$p.html()); },
            error:function(ev){
                alert('error ' + ev.target.status+' - '+ev.target.statusText);
            },
            success:function(data){
                if (data.error) {
                    alert(data.error)
                }
                else {
                    $('#'+progress).css("display", "none");
                    $('#'+file_id+'_ok').css("display", "block");
                    $('#'+file_id+'_ok').html('File successfully downloaded');
                }
            }
        });
    up.send();
}

function send_video (file_id, progress, type) {
    var
        $f = $('#'+file_id),
        $p = $('#'+progress),
        up = new uploader($f.get(0), {
            url:'ajax/upload.php',
            prefix:'file',
            type:type,
            progress:function(ev){ $p.html(((ev.loaded/ev.total)*100)+'%'); $p.css('width',$p.html()); },
            error:function(ev){
                alert('error ' + ev.target.status+' - '+ev.target.statusText);
            },
            success:function(data){
                if (data.error) {
                    alert(data.error)
                }
                else {
                    $('#'+progress).css("display", "none");
                    $('#'+file_id+'_ok').css("display", "block");
                    $('#'+file_id+'_ok').html('File successfully downloaded');
                    $('#video').css("display", "block");
                }
            }
        });
    up.send();
}

function show_profile (user_id) {
    $.post( 'ajax/profile.php', {
        'user_id' : user_id
    }, function (data) {
        $("#profile_abuses").html(data.abuses);
        $("#profile_reg_time").html(data.reg_time);
        $("#profile_window").css("display", "block");
        $("#profile_window").center();
        $("#reloadbtn").html('<button onclick="reload_photo('+user_id+', \'profile_photo\');" class="btn">reload photo</button>');
        $('#profile_photo').attr('src', '');
        reload_photo(user_id, 'profile_photo');
    }, 'JSON' );
}

function reload_photo(user_id, face_id) {
    $.post( 'ajax/new_photo.php', {
        'user_id' : user_id
    }, function (data) {
        $('#'+face_id).attr('src', ''+data.face+'');
    }, "json" );
}

jQuery.fn.center = function () {
    this.css("position","absolute");
    this.css("top", Math.max(0, (($(window).height() - $(this).outerHeight()) / 2) +
        $(window).scrollTop()) + "px");
    this.css("left", Math.max(0, (($(window).width() - $(this).outerWidth()) / 2) +
        $(window).scrollLeft()) + "px");
    return this;
}
function decrypt_comment(id, type) {

    console.log('decrypt_comment');

    var key = $("#key").text();
    var pass = $("#password").text();
    if (key.indexOf('RSA PRIVATE KEY')!=-1)
        pass = '';
    var e_text = $("#encrypt_comment_"+id).val();
    console.log('key='+key);
    console.log('pass='+pass);
    console.log('e_text='+e_text);

    if (pass) {
        text = atob(key.replace(/\n|\r/g,""));
        var decrypt_PEM = mcrypt.Decrypt(text, "\u0098nq\u0001\u009f\u00c9\u00d1\u00eb\u0012\u008dj\u000e\u00e0\u009d\u008f", hex_md5(pass), 'rijndael-128', 'ecb');
    }
    else
        decrypt_PEM = key;
    console.log('decrypt_PEM='+decrypt_PEM);

    var rsa2 = new RSAKey();
    rsa2.readPrivateKeyFromPEMString(decrypt_PEM); // N,E,D,P,Q,DP,DQ,C

    var decrypt_comment_ = rsa2.decrypt(e_text);

    $.post( 'ajax/save_decrypt_comment.php', {
        'id' : id,
        'comment' : decrypt_comment_,
        'type' : type
    }, function (data) {
        $("#comment_"+id).html(data);
    } );
}


function decrypt_comment_01 (id, type, miner_id, mcrypt_iv) {

    var key = $("#key").text();
    var pass = $("#password").text();
    if (key.indexOf('RSA PRIVATE KEY')!=-1)
        pass = '';
    var e_text = $("#encrypt_comment_"+id).val();

    if (miner_id > 0) { // если майнер, то коммент зашифрован нодовским ключем и тут его не расшифровать
        var comment = e_text;
    }
    else {
        if (pass) {
            text = atob(key.replace(/\n|\r/g, ""));
            var decrypt_PEM = mcrypt.Decrypt(text, mcrypt_iv, hex_md5(pass), 'rijndael-128', 'ecb');
        }
        else {
            decrypt_PEM = key;
        }
        var rsa2 = new RSAKey();
        rsa2.readPrivateKeyFromPEMString(decrypt_PEM); // N,E,D,P,Q,DP,DQ,C

        var comment = rsa2.decrypt(e_text);
    }
    // decrypt_comment может содержать зловред
    $.post( 'ajax/save_decrypt_comment.php', {
        'id' : id,
        'comment' : comment,
        'type' : type
    }, function (data) {
        console.log(data);
        $(".comment_"+id).html(data);
        console.log(".comment_"+id);
    }, 'HTML' );

}

function decrypt_message(id, type) {
    var key = $("#key").text();
    text = atob(key.replace(/\n|\r/g,""));
    var pass = $("#password").text();
    var e_text = $("#encrypt_comment_"+id).val();
    var decrypt_PEM = mcrypt.Decrypt(text, "\u0098nq\u0001\u009f\u00c9\u00d1\u00eb\u0012\u008dj\u000e\u00e0\u009d\u008f", pass, 'rijndael-128', 'ecb');

    var rsa2 = new RSAKey();
    rsa2.readPrivateKeyFromPEMString(decrypt_PEM); // N,E,D,P,Q,DP,DQ,C

    decrypt_comment = rsa2.decrypt(e_text);

    $.post( 'ajax/save_decrypt_comment.php', {
        'id' : id,
        'comment' : decrypt_comment,
        'type' : type
    }, function (data) {
        $("#comment_"+id).html(data);
    } );
}

function decrypt_admin_message(id) {
    var key = $("#key").text();
    text = atob(key.replace(/\n|\r/g,""));
    var pass = $("#password").text();
    var e_text = $("#encrypt_message_"+id).val();
    var decrypt_PEM = mcrypt.Decrypt(text, "\u0098nq\u0001\u009f\u00c9\u00d1\u00eb\u0012\u008dj\u000e\u00e0\u009d\u008f", pass, 'rijndael-128', 'ecb');

    var rsa2 = new RSAKey();
    rsa2.readPrivateKeyFromPEMString(decrypt_PEM); // N,E,D,P,Q,DP,DQ,C

    decrypt_comment = rsa2.decrypt(e_text);

    $.post( 'ajax/save_admin_decrypt_message.php', {
        'id' : id,
        'message' : decrypt_comment
    }, function (data) {

    } );
}