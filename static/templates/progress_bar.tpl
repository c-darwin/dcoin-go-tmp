<style>
	.progress_status {width:800px;}
	#progress_bar {display: block}


	@media  (max-width: 319px) {
		#progress_bar {display: none}
	}
	@media  (min-width: 319px) {
		.progress_status {width:150px;}
	}
	@media  (min-width: 410px) {
		.progress_status {width:240px;}
	}
	@media  (min-width: 510px) {
		.progress_status {width:340px;}
	}
	@media  (min-width: 600px) {
		.progress_status {width:200px;}
	}
	@media  (min-width: 768px) {
		.progress_status {width:180px;}
	}
	@media  (min-width: 810px) {
		.progress_status {width:220px;}
	}
	@media  (min-width: 900px) {
		.progress_status {width:280px;}
	}
	@media  (min-width: 1000px) {
		.progress_status {width:400px;}
	}
	@media  (min-width: 1100px) {
		.progress_status {width:500px;}
	}
	@media  (min-width: 1200px) {
		.progress_status {width:600px;}
	}
	@media  (min-width: 1300px) {
		.progress_status {width:700px;}
	}
	@media  (min-width: 1400px) {
		.progress_status {width:800px;}
	}
	@media  (min-width: 1500px) {
		.progress_status {width:900px;}
	}
</style>

<div style="float: left; margin-right: 15px; margin-top: 15px" id="progress_status_text"><?php echo $lng['your_status']?>: </div>
	<div style="float: left; position: relative">
		<div style=" display:table;" class="progress_status" onclick="dcNav({'target':{'hash':'#progress'}})">
			<div style="display: table-row;min-width:100%; font-size: 12px">
				<div style="display: table-cell;width: 33%;  "><?php echo $lng['progress_user']?></div>
				<div style="display: table-cell;width: 33%;  text-align: center "><?php echo $lng['progress_miner']?></div>
				<div style="display: table-cell;width: 33%;   text-align: right"><?php echo $lng['progress_boss']?></div>
			</div>
		</div>
		<div class="progress progress_status" style=" margin-top: 0px; margin-bottom:0px;" onclick="dcNav({'target':{'hash':'#progress'}})">
			<div id="progress_pct" class="progress-bar progress-bar-success" role="progressbar" aria-valuenow="10" aria-valuemin="0" aria-valuemax="100" style="width: <?php echo $tpl['progress_pct']?>%;">
				<?php echo $tpl['progress_pct']?>%
			</div>
		</div>
		<button type="button" class="close"  style="position: absolute; top: -4px; right: -15px" onclick="dcNav({'target':{'hash':'#interface/show_progress_bar=0'}})">×</button>
	</div>