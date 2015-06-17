
<script>

$('#send_to_net').bind('click', function () {

	$.post( 'ajax?controllerName=saveQueue', {
			'type' : '<?php echo $tpl['data']['type']?>',
			'time' : '<?php echo $tpl['data']['time']?>',
			'user_id' : '<?php echo $tpl['data']['user_id']?>',
			'funding_id' : <?php echo $tpl['del_id']?>,
			'signature1': $('#signature1').val(),
			'signature2': $('#signature2').val(),
			'signature3': $('#signature3').val()
			}, function (data) {
				//alert(data);
				fc_navigate ('wallets_list', {'alert': '<?php echo $lng['sent_to_the_net'] ?>'} );
			} );
} );

$("#main_div select").addClass( "form-control" );
$("#main_div input").addClass( "form-control" );
$("#main_div textarea").addClass( "form-control" );
$("#main_div button").addClass( "btn-outline btn-primary" );

</script>

<div id="main_div">
  <h1 class="page-header"><?php echo $lng['del_cf_funding_title']?></h1>

	<?php require_once( ABSPATH . 'templates/alert_success.php' );?>

    <div id="sign">
	
		<label><?php echo $lng['data']?></label>
		<textarea id="for-signature" style="width:500px;" rows="4"><?php echo "{$tpl['data']['type_id']},{$tpl['data']['time']},{$tpl['data']['user_id']},{$tpl['del_id']}"; ?></textarea>
	    <?php
		for ($i=1; $i<=$count_sign; $i++) {
			echo "<label>{$lng['sign']} ".(($i>1)?$i:'')."</label><textarea id=\"signature{$i}\" style=\"width:500px;\" rows=\"4\"></textarea>";
	    }
	    ?>
	    <br>
		<button class="btn" id="send_to_net"><?php echo $lng['send_to_net']?></button>

    </div>

	<input type="hidden" id="user_id" value="<?php echo $_SESSION['user_id']?>">
	<input type="hidden" id="time" value="<?php echo time()?>">
	<script>
		doSign();
		<?php echo !defined('SHOW_SIGN_DATA')?'$("#send_to_net").trigger("click");':'' ?>
	</script>
</div>