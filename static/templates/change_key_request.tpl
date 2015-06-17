<script>

$('#save').bind('click', function () {

	<?php echo !defined('SHOW_SIGN_DATA')?'':'$("#main").css("display", "none");	$("#sign").css("display", "block");' ?>
	$("#for-signature").val( '<?php echo "{$tpl['data']['type_id']},{$tpl['data']['time']},{$tpl['data']['user_id']}"; ?>,'+$("#to_user_id").val() );
	doSign();
	<?php echo !defined('SHOW_SIGN_DATA')?'$("#send_to_net").trigger("click");':'' ?>

});

$('#send_to_net').bind('click', function () {
	$.post( 'ajax?controllerName=saveQueue', {
			'type' : '<?php echo $tpl['data']['type']?>',
			'time' : '<?php echo $tpl['data']['time']?>',
			'user_id' : '<?php echo $tpl['data']['user_id']?>',
			'to_user_id' : $('#to_user_id').val(),
			'signature1': $('#signature1').val(),
			'signature2': $('#signature2').val(),
			'signature3': $('#signature3').val()
		}, function (data) {
			fc_navigate ('restoring_access', {'alert': '<?php echo $lng['sent_to_the_net'] ?>'} );
		}
	);

} );

</script>

	<h1 class="page-header"><?php echo $lng['request_access_to_the_account']?></h1>

	<?php require_once( ABSPATH . 'templates/alert_success.php' );?>
	
	<div id="main">
		<input type="text" id="to_user_id" class="form-control"><br>
		<button id="save" class="btn btn-outline btn-primary"><?php echo $lng['send_to_net'] ?></button>
	</div>

	<?php require_once( 'signatures.tpl' );?>

