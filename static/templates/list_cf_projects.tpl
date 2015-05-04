
<script>

$('#save').bind('click', function () {

	$("#change_host").css("display", "none");
	$("#sign").css("display", "block");
	$("#for-signature").val( '<?php echo "{$tpl['data']['type_id']},{$tpl['data']['time']},{$tpl['data']['user_id']}"; ?>,'+$("#currency_id").val()+','+$("#amount").val()+','+$("#end_time").val()+','+$("#latitude").val()+','+$("#longitude").val()+','+$("#category_id").val());
	doSign();
	<?php echo !defined('SHOW_SIGN_DATA')?'$("#send_to_net").trigger("click");':'' ?>
});

$('#send_to_net').bind('click', function () {

	$.post( 'ajax/save_queue.php', {
			'type' : '<?php echo $tpl['data']['type']?>',
			'time' : '<?php echo $tpl['data']['time']?>',
			'user_id' : '<?php echo $tpl['data']['user_id']?>',
			'currency_id' : $('#currency_id').val(),
			'amount' : $('#amount').val(),
			'end_time' : $('#end_time').val(),
			'latitude' : $('#latitude').val(),
			'longitude' : $('#longitude').val(),
			'category_id' : $('#category_id').val(),
			'signature1': $('#signature1').val(),
			'signature2': $('#signature2').val(),
			'signature3': $('#signature3').val()
		}, function (data) {
				fc_navigate ('new_cf_project', {'alert': '<?php echo $lng['sent_to_the_net'] ?>'} );
			}
	);
} );

</script>

	<h1 class="page-header"><?php echo $lng['new_cf_project_title']?></h1>

	<?php require_once( ABSPATH . 'templates/alert_success.php' );?>
	
	<div id="change_host">

		<form>
			<fieldset>
				<input type="text" placeholder="currency_id" id="currency_id" value=""><br>
				<input type="text" placeholder="amount" id="amount" value=""><br>
				<input type="text" placeholder="end_time" id="end_time" value="<?php echo time()+3600*24*10?>"><br>
				<input type="text" placeholder="latitude" id="latitude" value=""><br>
				<input type="text" placeholder="longitude" id="longitude" value=""><br>
				<input type="text" placeholder="category_id" id="category_id" value=""><br>
				<button type="submit" class="btn" id="save"><?php echo $lng['next']?></button>
			</fieldset>
		</form>

		<p><span class="label label-important"><?php echo $lng['limits']?></span> <?php echo $tpl['limits_text']?></p>

	</div>

	<?php require_once( 'signatures.tpl' );?>

