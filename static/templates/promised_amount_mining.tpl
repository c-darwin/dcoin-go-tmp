
<script>

$('#send_to_net').bind('click', function () {
	if ( $('#amount').val() > 0 ) {
		$.post( 'ajax?controllerName=saveQueue', {
				'type' : '<?php echo $tpl['data']['type']?>',
				'time' : '<?php echo $tpl['data']['time']?>',
				'user_id' : '<?php echo $tpl['data']['user_id']?>',
				'promised_amount_id' :  $('#promised_amount_id').val(),
				'amount' :  $('#amount').val(),
				'signature1': $('#signature1').val(),
				'signature2': $('#signature2').val(),
				'signature3': $('#signature3').val()
			}, function(data){
			fc_navigate ('promised_amount_list', {'alert': '<?php echo $lng['sent_to_the_net'] ?>'} );
		});
	}
	else	{
		alert('null amount');
	}
} );

$("#main_div textarea").addClass( "form-control" );

</script>
<div id="main_div">
<h1 class="page-header"><?php echo $lng['mining']?></h1>
<ol class="breadcrumb">
	<li><a href="#mining_menu"><?php echo $lng['mining'] ?></a></li>
	<li><a href="#promised_amount_list"><?php echo $lng['promised_amount_title'] ?></a></li>
	<li class="active"><?php echo $lng['mining'] ?></li>
</ol>

    <div id="sign_banknote">
	
		<label><?php echo $lng['data']?></label>
		<textarea id="for-signature" style="width:500px;" rows="4"><?php echo "{$tpl['data']['type_id']},{$tpl['data']['time']},{$tpl['data']['user_id']},{$tpl['promised_amount_id']},{$tpl['amount']}"?></textarea>
	    <?php
	for ($i=1; $i<=$count_sign; $i++) {
		echo "<label>{$lng['sign']} ".(($i>1)?$i:'')."</label><textarea id=\"signature{$i}\" style=\"width:500px;\" rows=\"4\"></textarea>";
	    }
	    ?>
		<br>
	    <button class="btn" id="send_to_net"><?php echo $lng['send_to_net']?></button>

    </div>

	<input type="hidden" id="amount" value="<?php echo $tpl['amount']?>">
	<input type="hidden" id="promised_amount_id" value="<?php echo $tpl['promised_amount_id']?>">
	<script>
		doSign();
		<?php echo !defined('SHOW_SIGN_DATA')?'$("#send_to_net").trigger("click");':'' ?>
	</script>
</div>