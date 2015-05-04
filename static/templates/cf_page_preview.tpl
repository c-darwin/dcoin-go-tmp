<?php
if (isset($_REQUEST['blurb_img'])) {
?>
<!DOCTYPE html>
<html lang="en">

<head>

	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<meta name="description" content="">
	<meta name="author" content="">

	<title>Preview</title>

	<!-- Bootstrap Core CSS -->
	<link href="css/bootstrap.min.css" rel="stylesheet">

	<!-- MetisMenu CSS -->
	<link href="css/plugins/metisMenu/metisMenu.min.css" rel="stylesheet">

	<!-- Custom CSS -->
	<link href="css/sb-admin.css" rel="stylesheet">

	<!-- Custom Fonts -->
	<link href="css/font-awesome.css" rel="stylesheet">

	<!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
	<!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
	<!--[if lt IE 9]>
	<script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
	<script src="https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
	<![endif]-->
<?php
}
?>
	<style>
		body {background-color:#F7FAFA; padding:0; margin:0}

		.left {
			float: left;
		}
		.bold {
			font-weight: bold;
		}
		ul {
			list-style: disc outside none;
		}
		.right {
			float: right;
		}
		.white {
			color: #FFF;
		}
	#page-wrapper{background-color:#F7FAFA;}
	#cf_active_menu a{border-color: #2BDE73; color: #000; border-bottom: 5px solid  #2BDE73;padding-bottom: 7px;}
	#cf_active_menu a:link{color: #000;}
	#cf_active_menu a:visited{color: #000;}
	#cf_active_menu a:hover{color: #000;}
	#cf_active_menu a:active{color: #000;}
		<?php
		if (!$user_id) {
			echo "#page-wrapper{margin:0}\n";
		}
		?>
	</style>
	<link href="<?php echo $tpl['cf_url']?>css/cf.css" rel="stylesheet">
	<link href="<?php echo $tpl['cf_url']?>css/social-buttons.css" rel="stylesheet">
<?php
if (1<0) {
?>
</head>

<body>
<!-- container -->
<?php
}
?>

<script>
	$('#send_comment').bind('click', function () {

		<?php echo !defined('SHOW_SIGN_DATA')?'':'$("#sign").css("display", "block"); $("#comment_div").css("display", "none");' ?>

		$("#for-signature").val( '<?php echo "{$tpl['comment_data']['type_id']},{$tpl['data']['time']},{$tpl['data']['user_id']},{$tpl['project_id']},{$tpl['lang_id']}"; ?>,'+$("#comment").val());

		doSign();

		<?php echo !defined('SHOW_SIGN_DATA')?'$("#send_to_net").trigger("click");':'' ?>

	});

	$('#comment_send_to_net').find('#send_to_net').bind('click', function () {
			$('#wrapper').spin();
			$.post( 'ajax/save_queue.php', {
					'type' : '<?php echo $tpl['comment_data']['type']?>',
					'time' : '<?php echo $tpl['data']['time']?>',
					'user_id' : '<?php echo $tpl['data']['user_id']?>',
					'project_id' : <?php echo $tpl['project_id']?>,
					'lang_id' :  <?php echo $tpl['lang_id']?>,
					'comment' : $('#comment').val(),
					'signature1': $('#signature1').val(),
					'signature2': $('#signature2').val(),
					'signature3': $('#signature3').val()
			}, function (data) {
				$("#main_comment_div").html( '<div class="alert alert-success">Если Ваш комментарий не содержит ошибок и не были превышены лимиты, тогда он отобразится на этой странице через нескольких минут</div>');
				$('#wrapper').spin(false);

			}
	);
	});

	$('#contribute_now').bind('click', function () {
		$("#project_info").css("display", "none");
		$("#payment_mode").css("display", "block");
	});

	$('#contribute_now_step1').bind('click', function () {
		if ($('input[name=mode]:checked').val()=='1') {
			$("#payment_mode").css("display", "none");
			$("#payment_method").css("display", "block");
		}
		else if ($('input[name=mode]:checked').val()=='2') {
			//javascript:location.href='<?php echo $tpl['config']['pool_url']?>';
			fc_navigate('cf_start');
		}
	});

	$('#contribute_now_step2').bind('click', function () {

			$('#wrapper').spin();
			$.post( '<?php echo $tpl['cf_url']?>ajax/available_coins.php', {
					'dc_currency_id' : '<?php echo $tpl['project']['currency_id']?>',
					'currency_id' : '1001',
					'amount' : $("#amount_usd_total").val()
				}, function (data) {
					if (data.error) {
						$("#payment_error").html( '<div class="alert alert-danger alert-dismissable"><button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>'+data.error+'</div>');
					}
					else if (data.success) {
						$("#payment_method").css("display", "none");
						if ($( "#p_method option:selected" ).val() == 'gspay') {
							$("#gspay_form").submit();
						}
						else if ($( "#p_method option:selected" ).val() == 'ik') {
							$("#ik_form").submit();
						}
						else if ($( "#p_method option:selected" ).val() == 'pm') {
							$("#pm_form").submit();
						}

					}
					$('#wrapper').spin(false);
				}, 'JSON'
			);

	});

	function change_amount() {
		var amount = $('#amount_usd').val();
		if ($( "#p_method option:selected" ).val() == 'gspay') {
			var total_usd = (amount * 0.933) - 1;
		}
		else if ($( "#p_method option:selected" ).val() == 'ik') {
			var total_usd = (amount * 0.99);
		}
		else if ($( "#p_method option:selected" ).val() == 'pm') {
			var total_usd = (amount * 0.99);
		}
		total_usd = total_usd.toFixed(2);
		if (total_usd<0)
			total_usd = 0;

		$('#pm_amount').val(amount);
		$('#ik_amount').val(amount);
		$('#gspay_amount').val(amount);
		$('#amount_usd_total').val(total_usd);
		$('#amount_dc_total').text('~ '+total_usd);
	}

	$('#amount_usd').bind('keyup', function(event) {
		change_amount();
	});

	$('#p_method').on('change', function(event) {
		change_amount();
	});



</script>

<!-- Page Content -->
<div class="container" style="margin-left: auto;margin-right: auto; margin-bottom:50px; margin-top:50px; width: 1000px;padding: 0px 0px 0px 0px ">
<!--
	<div class="width_max">

		<ul class="nav navbar-nav navbar-left" style="padding-top:10px">
			<button type="button" class="btn btn-outline btn-default">Explore</button>
			<button type="button" class="btn btn-outline btn-default">Start your campaign</button>
		</ul>

		<ul class="nav navbar-nav navbar-right">
			<li class="dropdown">
				<a href="#" class="dropdown-toggle" data-toggle="dropdown">Language <span class="caret"></span></a>
				<ul class="dropdown-menu" role="menu">
					<li><a href="#">Русский</a></li>
					<li><a href="#">English</a></li>
				</ul>
			</li>
		</ul>

	</div>
-->

	<img src="<?php echo $tpl['head_img']?>" style="width:1000px; height:150px">

	<div class="menu width_max" style="height: 42px">

		<ul class="list-inline left bold">
			<?php
			foreach ($tpl['pages_array'] as $page_type) {
				$counter = '';
				if ($page_type=='funders')
					$counter = ' <span class="count h6 bg-grey-dark white">'.$tpl['project']['count_funders'].'</span>';
				else if ($page_type=='comments')
					$counter = ' <span class="count h6 bg-grey-dark white">'.$tpl['project']['count_comments'].'</span>';
				echo "<li ";
				echo ($tpl['page']==$page_type)?'id="cf_active_menu">':'>';
				if (!$user_id)
					echo "<a href='?id-{$tpl['project_id']}-{$tpl['lang_id']}-{$page_type}'>";
				else
					echo "<a href=\"#\" onclick=\"fc_navigate('cf_page_preview', {'only_project_id':{$tpl['project_id']}, 'lang_id':{$tpl['lang_id']}, 'page':'{$page_type}'})\">";
				echo "{$lng['cf_'.$page_type]}{$counter}</a></li>";
			}
			?>
		</ul>

		<ul class="list-inline right bold" style="margin-right:10px">

			<li><?php echo $tpl['project']['project_currency_name']?></li>
			<li>Project ID: <?php echo $tpl['project']['id']?></li>
			<?php echo $tpl['project']['country']?'<li><i class="fa  fa-map-marker  fa-fw"></i>  '.$tpl['project']['country'].', '.$tpl['project']['city'].'</li>':'' ?>
			<li>
				<?php
				if ($user_id)
					echo "<a href=\"#\" onclick=\"fc_navigate('cf_catalog', {'category_id':{$tpl['project']['category_id']}})\">";
				else
					echo "<a href='?category-{$tpl['project']['category_id']}'>";
				?>

				<i class="fa  fa-folder-open-o  fa-fw"></i> <?php echo $lng['cf_category'][$tpl['project']['category_id']]?></a></li>
		</ul>
	</div>

	<!-- /.row -->
	<div class="clearfix"></div>

	<div class="well" style="background-color:#fff;margin:auto; width:1000px; padding-top:0px">

		<div class="row">
			<div id="project_info">
			<ul class="list-inline lng" style="margin-left:20px; margin-top:13px">
				<?php
				if ($tpl['page']!='funders')
				foreach ($tpl['project']['lang'] as $data_id=>$lang_id) {
					$num = '';
					if ($tpl['page'] == 'comments')
						$num = ' <span class="h6" style="color: #000; border-radius:3px;background:#ddd;font-weight:normal;padding:2px 5px; font-size: 13px;">'.(int)$tpl['project']['lang_comments'][$lang_id].'</span>';
					if ($tpl['lang_id']!=$lang_id) {
						if ($user_id)
							echo "<li><a href=\"#\" onclick=\"fc_navigate('cf_page_preview', {'only_project_id':{$tpl['project_id']}, 'lang_id':{$lang_id}, 'page':'{$tpl['page']}'})\">{$tpl['cf_lng'][$lang_id]}</a>{$num}</li> ";
						else
							echo "<li><a href='?id-{$tpl['project_id']}-{$lang_id}-{$tpl['page']}'>{$tpl['cf_lng'][$lang_id]}</a>{$num}</li> ";
					}
					else {
						if (sizeof($tpl['project']['lang'])==1)
							echo "";
						else
							echo "<li>{$tpl['cf_lng'][$lang_id]}{$num}</li> ";
					}
				}
				?>
			</ul>

			<div style="width:620px; float:left; margin: 5px 35px 0px 25px;">
			<?php
			if ($tpl['page'] == 'home') {
				$project_url = $tpl['cf_url'].'?'.$tpl['project']['project_currency_name'].'-'.$tpl['lang_id'];
				if ($tpl['video_url_id'])
					echo '<iframe width="620" height="413" src="http://www.youtube.com/embed/'.$tpl['video_url_id'].'" frameborder="0" allowfullscreen></iframe>';
				else
				 echo '<img src="'.$tpl['picture'].'?r='.rand().'" width="620" height="413">';
				if ($tpl['cf_url'])
					echo '<div class="text-center" style="margin: 5px 0px 5px 0px"><a class="btn btn-social-icon btn-twitter" href="http://twitter.com/intent/tweet?text='.$project_url.'" target="_blank"><i class="fa fa-twitter"></i></a> <a class="btn btn-social-icon btn-vk" href="http://vk.com/share.php?url='.$project_url.'" target="_blank"><i class="fa fa-vk"></i></a> <a href="https://www.facebook.com/sharer/sharer.php?u='.$project_url.'" target="_blank" class="btn btn-social-icon btn-facebook"><i class="fa fa-facebook"></i></a> <a  href="https://plus.google.com/share?url='.$project_url.'" target="_blank" class="btn btn-social-icon btn-google-plus"><i class="fa fa-google-plus"></i></a></div>';
			}
			else if ($tpl['page'] == 'news')
				echo '<img src="'.$tpl['news_img'].'?r='.rand().'" width="620">';
			else if ($tpl['page'] == 'funders') {
				foreach ($tpl['funders'] as  $data) {
					echo '<div style="overflow: hidden;padding: 15px 15px 15px 0;border-bottom: 1px solid #D9D9DE;"><img src="'.$data['avatar'].'" style="width: 80px; height: 80px; float:left; margin: 0 15px"><div><p><strong>'.$data['name'].'</strong></p><p>'.$data['time'].'</p></div></div>';
				}
			}
			else if ($tpl['page'] == 'comments') {

				echo '<div id="main_comment_div">';
				if (!$tpl['project']['funder'] && $tpl['project']['user_id']!=$user_id)
					echo '<div class="alert alert-info">'.$lng['comments_only_for_funders'].'</div>';
				else {
					echo '<div id="comment_div"><div class="alert alert-info"><strong>'.$lng['limits'].':</strong> '.$lng['comments_limits'].'</div><div><textarea id="comment" class="form-control" rows="3" maxlength="140"></textarea></div><div><button type="button" class="btn btn-outline btn-primary btn-lg btn-block" style="margin-bottom: 20px; margin-top: 5px" id="send_comment">'.$lng['send'].'</button></div></div>';
					echo '<div style="margin-bottom: 20px" id="comment_send_to_net">';
					require_once( 'signatures.tpl' );
					echo '</div>';
				}
				echo '</div>';

				foreach ($tpl['comments'] as  $data) {
					$bd = '';
					if ($tpl['project']['user_id']==$data['user_id']) {
						$bd = 'background-color:#E8F6FF;';
						$data['name'] = $data['name'].' <span style="color: #ff0000">(Creator)</span> ';
					}
					echo '<div style="overflow: auto;padding: 15px 15px 15px 0;border-bottom: 1px solid #D9D9DE; '.$bd.' "><img src="'.$data['avatar'].'" style="width: 80px; height: 80px; float:left; margin: 0 15px"><div style="overflow: auto;"><p><strong>'.$data['name'].'</strong> <span style="color: #999">'.$data['time'].'</span></p><p>'.$data['comment'].'</p></div></div>';
				}
			}
			else if ($tpl['page'] == 'payment') {
				echo '<div style="overflow: hidden; padding: 15px 15px 15px 0;border-bottom: 1px solid #D9D9DE;">'.str_replace('[currency_name]', 'D'.$tpl['project']['currency'], $lng['payment_credited']).'</div>';
			}
			?>
			</div>


				<div id="project-info" style="overflow:auto;">
					<div style="margin-left:18px">
						<h1><?php echo $tpl['project']['funding']?></h1>
						<p><?php echo $lng['cf_page_preview_pledged_of']?> <?php echo $tpl['project']['amount']?> D<?php echo $tpl['project']['currency']?> <?php echo $lng['cf_page_preview_goal']?> </p>
						<h1><?php echo $tpl['project']['days']?></h1>
						<p><?php echo $lng['days_to_go']?></p>
						<p style="font-weight: normal"><?php echo $lng['start_date']?> <?php echo $tpl['project']['start_date']?></p>

						<?php
						if (@$tpl['project']['ended']!=1)
						{
							if ($user_id)
								echo "<button type=\"button\" class=\"btn btn-success\" style=\"width:240px; height:50px\" onclick=\"fc_navigate('wallets_list', {'project_id':{$tpl['project']['id']}})\"><strong>".($lng['contribute_now'])."</strong></button>";
							else if (!empty($tpl['project']['ps']))
								echo "<button type=\"button\" class=\"btn btn-success\" style=\"width:240px; height:50px\" id=\"contribute_now\"><strong>".($lng['contribute_now'])."</strong></button>";
							else
								echo "<button type=\"button\" class=\"btn btn-success\" style=\"width:240px; height:50px\" onclick=\"fc_navigate('cf_start')\"><strong>".($lng['contribute_now'])."</strong></button>";
						}
						?>
					</div>

					<div class="well" style="background-color:#E8F6FF; border:0px; pading:10px; margin-top:25px; width:280px; height:140px">

						<div style="width: 100px; float: left;margin-right:10px"><img src="<?php echo $tpl['project']['author']['avatar']?>" style="width:100px; height: 100px"></div>

						<div>
							<h4 style="margin-top:0px"><?php echo $tpl['project']['author']['name']?></h4>
							<h5><?php echo $tpl['project']['author']['created']?> <?php echo $lng['created']?><br><?php echo $tpl['project']['author']['backed']?>  <?php echo $lng['backed']?></h5>
							<div class="clearfix"></div>
						</div>
					</div>
					<div class="clearfix"></div>

				</div>

				<?php

				if ($tpl['page'] == 'home') {
					echo '<img src="'.$tpl['description_img'].'?r='.rand().'" style="width:990px; margin:auto" '.($tpl['links']?'usemap="#Navigation"':'').'>';
					if ($tpl['links']) {
						echo '<map name="Navigation">';
						foreach ($tpl['links'] as $data)
							echo "<area shape=\"rect\" coords=\"{$data[1]},{$data[2]},{$data[3]},{$data[4]}\" href=\"{$data[0]}\" target='_blank'>";
						echo '</map>';
					}
				}


				?>

			</div>

			<div id="payment_mode" style="display: none;margin-top:35px;margin-left:15px;margin-right:15px;"">

				<div class="form-horizontal">
					<fieldset>
						<div class="form-group">
							<label class="col-md-4 control-label" for="radios"><?php echo $lng['select_the_mode']?></label>
							<div class="col-md-4">
								<div class="radio">
									<label>
										<input name="mode" value="1" checked="checked" type="radio">
										<?php echo $lng['cf_send_money_wo_sign_up']?>
									</label>
								</div>
								<div class="radio">
									<label>
										<input name="mode" value="2" type="radio">
										<?php echo $lng['cf_send_money_with_sign_up']?>
									</label>
								</div>
							</div>
						</div>
						<div class="form-group">
							<label class="col-md-4 control-label" for="singlebutton"></label>
							<div class="col-md-4">
								<button id="contribute_now_step1" class="btn btn-outline btn-primary"><?php echo $lng['next']?></button>
							</div>
						</div>

					</fieldset>
				</div>
			</div>

			<div id="payment_method" style="display: none;margin-top: 35px;margin-left:15px;margin-right:15px;">
				<div id="payment_error"></div>

				<div class="form-horizontal" method="post">
					<fieldset>
						<div class="form-group">
							<label class="col-md-4 control-label" for="selectbasic">Method</label>
							<div class="col-md-4">
								<select id="p_method" name="p_method" class="form-control">
									<?php
									foreach ($tpl['config']['cf_ps'] as $ps_id=>$data) {
										if (!empty($tpl['project']['ps']['ps'.$ps_id]))
											echo "<option value='{$data[0]}'>{$data[1]}</option>";
									}
									?>
								</select>
							</div>
						</div>

						<div class="form-group">
							<label class="col-md-4 control-label" for="amount_usd">Amount</label>
							<div class="col-md-4">
								<div class="input-group">
									<input id="amount_usd" name="amount_usd" class="form-control" type="text">
									<span class="input-group-addon">USD</span>
								</div>
							</div>
						</div>
						<div class="form-group">
							<label class="col-md-4 control-label" for="textinput">D<?php echo $tpl['project']['currency']?></label>
							<div class="col-md-4">
								<p class="form-control-static" id="amount_dc_total">0</p>
								<input type="hidden" id="amount_usd_total">
							</div>
						</div>
						<div class="form-group">
							<label class="col-md-4 control-label" for="singlebutton"></label>
							<div class="col-md-4">
								<button id="contribute_now_step2" class="btn btn-outline btn-primary"><?php echo $lng['next']?></button>
								<button id="p_submit" class="btn btn-outline btn-primary" style="display: none"><?php echo $lng['next']?></button>
							</div>
						</div>
					</fieldset>
				</div>
			</div>


		</div>
		<!-- /.row -->

		<form id="ik_form" method="post" action="https://sci.interkassa.com/" enctype="utf-8" style="display: none">
			<input type="hidden" name="ik_co_id" value="5400d4e9bf4efcf8492c867c" />
			<input type="hidden" name="ik_pm_no" value="ID_4233" />
			<input type="hidden" name="ik_cur" value="USD" />
			<input type="hidden" name="ik_suc_u" value="http://dcrowd.org/?id-<?php echo $tpl['project_id']?>-<?php echo $tpl['lang_id']?>-payment"/>
			<input type="hidden" name="ik_suc_m" value="get" />
			<input type="hidden" name="ik_fal_u" value="http://dcrowd.org/?id-<?php echo $tpl['project_id']?>-<?php echo $tpl['lang_id']?>">
			<input type="hidden" name="ik_fal_m" value="get" />
			<input type="hidden" name="ik_desc" value="cf-<?php echo $tpl['project']['id']?>" />
			<input id="ik_amount" name="ik_am" type="text">
			<input type="submit">
		</form>

		<form id="pm_form" method="post" action="https://perfectmoney.is/api/step1.asp" style="display: none">
			<input type="hidden" name="PAYEE_ACCOUNT" value="U8777532">
			<input type="hidden" name="PAYEE_NAME" value="DcoinSimple">
			<input type="hidden" name="PAYMENT_ID" value="cf-<?php echo $tpl['project']['id']?>">
			<input type="hidden" name="PAYMENT_UNITS" value="USD">
			<input type="hidden" name="STATUS_URL" value="http://DcoinSimple.com/pm.php">
			<input type="hidden" name="PAYMENT_URL" value="http://dcrowd.org/?id-<?php echo $tpl['project_id']?>-<?php echo $tpl['lang_id']?>-payment">
			<input type="hidden" name="PAYMENT_URL_METHOD" value="LINK">
			<input type="hidden" name="NOPAYMENT_URL" value="http://dcrowd.org/?id-<?php echo $tpl['project_id']?>-<?php echo $tpl['lang_id']?>">
			<input type="hidden" name="NOPAYMENT_URL_METHOD" value="LINK">
			<input type="hidden" name="SUGGESTED_MEMO" value="Dcoins">
			<input type="hidden" name="BAGGAGE_FIELDS" value="">
			<input id="pm_amount" name="PAYMENT_AMOUNT" type="text">
			<input type="submit" name="PAYMENT_METHOD">
		</form>

		<form id="gspay_form" method=post action="https://secure.redirect2pay.com/payment/pay.php" style="display: none">
			<input type=hidden name="siteID" value="117618">
			<input type=hidden name='OrderDescription[1]' value='cf-<?php echo $tpl['project']['id']?>'>
			<input id="gspay_amount"  type=hidden name='Amount[1]' value='0'>
			<input type="hidden" name="OrderID" value="cf-<?php echo $tpl['project']['id']?>" />
			<input type="hidden" name="returnURL"  value="http://dcrowd.org/?id-<?php echo $tpl['project_id']?>-<?php echo $tpl['lang_id']?>" />
			<input type="hidden" name="ApproveURL" value="http://dcrowd.org/?id-<?php echo $tpl['project_id']?>-<?php echo $tpl['lang_id']?>-payment" />
			<input type="hidden" name="DeclineURL"  value="http://dcrowd.org/?id-<?php echo $tpl['project_id']?>-<?php echo $tpl['lang_id']?>" />
			<input type="submit" value="Checkout">
		</form>


	</div>
</div>

<?php
if (1<0) {
?>

<!-- /.container -->

<!-- jQuery Version 1.11.0 -->
<script src="js/jquery-1.11.0.js"></script>

<!-- Bootstrap Core JavaScript -->
<script src="js/bootstrap.min.js"></script>


</body>

</html>

<?php
}
?>