<!DOCTYPE html>
<html lang="en">

<head>

<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="description" content="">
<meta name="author" content="">

<title>Decentralized CrowdFunding for any projects</title>

<!-- Bootstrap Core CSS -->
<link href="<?php echo $tpl['cf_url']?>css/bootstrap.min.css" rel="stylesheet">

<!-- MetisMenu CSS -->
<link href="<?php echo $tpl['cf_url']?>css/plugins/metisMenu/metisMenu.min.css" rel="stylesheet">

<!-- Custom CSS -->
<link href="<?php echo $tpl['cf_url']?>css/sb-admin.css" rel="stylesheet">

<!-- Custom Fonts -->
<link href="<?php echo $tpl['cf_url']?>css/font-awesome.css" rel="stylesheet">

<!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
<!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
<!--[if lt IE 9]>
<script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
<script src="https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
<![endif]-->
<script src="<?php echo $tpl['cf_url']?>js/index.js"></script>
<script src="<?php echo $tpl['cf_url']?>js/jquery-1.11.0.js"></script>
<link rel="stylesheet" media="all" type="text/css" href="<?php echo $tpl['cf_url']?>css/jquery-ui.css" />
	<style>
		#page-wrapper{
			margin: 0px 10% 0px 10%;
			border: 0;
			min-height: 550px;
		}
		#wrapper{height: 100%;}
		#dc_content{
			height: 550px;
			vertical-align: middle;

	</style>

</head>

<body>

<div id="wrapper">

	<nav class="navbar navbar-default navbar-fixed-top" role="navigation" style="margin-bottom: 0">
		<div class="navbar-header">
			<button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".sidebar-collapse">
				<span class="sr-only">Toggle navigation</span>
				<span class="icon-bar"></span>
				<span class="icon-bar"></span>
				<span class="icon-bar"></span>
			</button>
			<a class="navbar-brand" href="<?php echo $tpl['cf_url']?>" style="display: block; /* or inline-block; I think IE would respect it since a link is an inline-element */
	                   background: url(<?php echo $tpl['cf_url']?>img/logo-small.png) center left no-repeat;
	                   text-align: left;
	                   background-size: 40px 40px;
	                   padding-left: 40px; margin-left: 15px; margin-right: 0px; line-height: 12px"><nobr>Dcoin <span style="font-size: 12px">v<?php echo $tpl['ver']?></span></nobr><br><span style="font-size: 12px">All the projects are taken from Dcoin blockchain</span></a>
		</div>
		<!-- /.navbar-header -->

		<ul class="nav navbar-top-links navbar-right">
			<li>
				<a href="#" onclick="fc_navigate('cf_start');return false;">Start Your Campaign</a>
			</li>
			<li>
				<a href="http://dcoin.me">About Dcoin</a>
			</li>
			<li>
				<a href="http://dcoinforum.org">Forum</a>
			</li>
			<li>
				<a href="http://en.dcoinwiki.com/">Wiki</a>
			</li>
			<li>
				<a href="mailto:hello@dcoin.me">Contact</a>
			</li>
			<!-- /.dropdown -->
			<li class="dropdown">
				<a class="dropdown-toggle" data-toggle="dropdown" href="#">
					<i class="fa  fa-globe fa-fw"></i> Language: <?php echo $tpl['cf_lang'][$lang]?> <i class="fa fa-caret-down"></i>
				</a>
				<ul class="dropdown-menu dropdown-user">
					<li><a href="#" onclick="fc_navigate('cf_catalog', 'lang=1'); load_menu();">English</a>
					</li>
					<li><a href="#" onclick="fc_navigate('cf_catalog', 'lang=42'); load_menu();">Русский</a>
					</li>
				</ul>
				<!-- /.dropdown-user -->
			</li>
			<!-- /.dropdown -->
		</ul>
		<!-- /.navbar-top-links -->

	</nav>




<div id="page-wrapper">
	<div class="row">
		<div class="col-lg-12">
			<div id="dc_content"></div>

		</div>
		<!-- /.col-lg-12 -->
	</div>
	<!-- /.row -->
</div>
<!-- /#page-wrapper -->

</div>
<!-- /#wrapper -->

	<script src="<?php echo $tpl['cf_url']?>js/bootstrap.min.js"></script>
  	<script>
    <?php
	echo "$( document ).ready(function() {{$tpl['nav']}});\n";
	?>
	</script>

	<script src="<?php echo $tpl['cf_url']?>js/markerclusterer.js"></script>

	<script type="text/javascript" src="<?php echo $tpl['cf_url']?>js/jquery-ui.min.js"></script>

	<script type="text/javascript" src="<?php echo $tpl['cf_url']?>js/jquery-ui-timepicker-addon.js"></script>
	<script type="text/javascript" src="<?php echo $tpl['cf_url']?>js/jquery-ui-sliderAccess.js"></script>

	<script language="JavaScript" type="text/javascript" src="<?php echo $tpl['cf_url']?>js/spin.js"></script>


<script>
	(function ($) {
		$.fn.spin = function (opts, color) {
			var presets = {
				"tiny": {
					lines: 8,
					length: 2,
					width: 2,
					radius: 3
				},
				"small": {
					lines: 8,
					length: 4,
					width: 3,
					radius: 5
				},
				"large": {
					lines: 10,
					length: 8,
					width: 4,
					radius: 8
				}
			};
			if (Spinner) {
				return this.each(function () {
					var $this = $(this),
							data = $this.data();

					if (data.spinner) {
						data.spinner.stop();
						delete data.spinner;
					}
					if (opts !== false) {
						if (typeof opts === "string") {
							if (opts in presets) {
								opts = presets[opts];
							} else {
								opts = {};
							}
							if (color) {
								opts.color = color;
							}
						}
						data.spinner = new Spinner($.extend({
							color: $this.css('color')
						}, opts)).spin(this);
					}
				});
			} else {
				throw "Spinner class not available.";
			}
		};
	})(jQuery);

</script>
<!-- Yandex.Metrika counter -->
<script type="text/javascript">
	(function (d, w, c) {
		(w[c] = w[c] || []).push(function() {
			try {
				w.yaCounter25998519 = new Ya.Metrika({id:25998519,
					webvisor:true,
					clickmap:true,
					trackLinks:true,
					accurateTrackBounce:true});
			} catch(e) { }
		});

		var n = d.getElementsByTagName("script")[0],
			s = d.createElement("script"),
			f = function () { n.parentNode.insertBefore(s, n); };
		s.type = "text/javascript";
		s.async = true;
		s.src = (d.location.protocol == "https:" ? "https:" : "http:") + "//mc.yandex.ru/metrika/watch.js";

		if (w.opera == "[object Opera]") {
			d.addEventListener("DOMContentLoaded", f, false);
		} else { f(); }
	})(document, window, "yandex_metrika_callbacks");
</script>
<noscript><div><img src="//mc.yandex.ru/watch/25998519" style="position:absolute; left:-9999px;" alt="" /></div></noscript>
<!-- /Yandex.Metrika counter -->
</body>
</html>
