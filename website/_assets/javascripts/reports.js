//= require libraries/jquery
//= require libraries/materialize.min
//= require libraries/lightbox.min
//= require libraries/charts

(function($){
	return{
		preloaderHtml : '<div class="reports-preloader"><div class="row"><div class="col m7 s12"><div class="preloader-section"><div class="section-head"><div class="head-animation"><span>lorem ipsum</span></div></div><div class="section-body"><div class="paragraphs"><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div></div></div></div><div class="preloader-section long"><div class="section-head"><div class="head-animation long"><span>lorem ipsum</span></div></div><div class="section-body"><div class="paragraphs"><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div></div></div></div></div><div class="col m5 s12"><div class="preloader-section long"><div class="section-head"><div class="head-animation"><span>lorem ipsum</span></div></div><div class="section-body"><div class="paragraphs"><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div></div></div></div><div class="preloader-section long"><div class="section-head"><div class="head-animation long"><span>lorem ipsum</span></div></div><div class="section-body"><div class="paragraphs"><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div><div class="line"></div></div></div></div></div></div></div>',

		init : function(){
			lightbox.option({
				'resizeDuration': 200,
				'wrapAround': true
			})
			this.reports();
		},

		reports : function(){
			this.reprtsScreenshots();
			this.loadFirstPage();
			this.loadReportPages();
		},

		loadFirstPage : function(){
			if($(".reports .reports-nav a.active").attr("href") === "/reports/overview.html"){ 
				$(".reports .reports-content").html(this.preloaderHtml);

				setTimeout(function(){
					$.ajax({
						method: "get",
						url : "/reports/overview.html",
						success : function(data){
							$(".reports .reports-content").html(data);
						}
					})
				}, 1);
			}
		},

		reprtsScreenshots : function(){
			if($(".reports").length && $(".screenshots").length){
				$(".screenshots .screenshots-slider ul.thumbs li").on("click", function(){
					if(!$(this).hasClass("active")){
						$(this).parent().find("li.active").removeClass("active");
						$(this).addClass("active");

						var index = $(this).index();
						$(this).parents(".screenshots-slider").find(".screenshots-images li.active").removeClass("active");
						$(this).parents(".screenshots-slider").find(".screenshots-images li").eq(index).addClass("active");
					}
				})
			}
		},

		loadReportPages : function(){
			if($(".reports").length){
				var $this = this;
				$(".reports .reports-nav a").on("click", function(e) {
					e.preventDefault();

					if(!$(this).hasClass("active") && $(this).attr("href").length){
						var _this = $(this);
						$(".reports .reports-content").html($this.preloaderHtml);

						$.ajax({
							method: "get",
							url : $(this).attr("href"),
							success : function(data){
								_this.parents("ul").find("a.active").removeClass("active");

								if(_this.parents(".sub-menu").length){
									_this.parents(".sub-menu").parent("li").children("a").addClass("active");
								}

								_this.addClass("active");

								setTimeout(function(){
									$(".reports .reports-content").html(data);
									$this.afterAjax();
								}, 1)
							},
							error : function(){
								alert("failed");
							}
						})
					}
				})
			}
		},

		afterAjax : function(){
			this.subitemTable();
			this.scorebox();
			this.collapseBar();
		},

		subitemTable : function(){
			if($("table.si-table").length){
				$("table.si-table tbody tr").on("click", function(){
					var id = $(this).attr("id");
					$(this).toggleClass("opened");

					$(this).parent("tbody").find("tr[id*='"+ id +"_']").toggleClass("active");
				})
			}
		},

		scorebox : function(){
			if($(".scorebox").length){
				$(".scorebox").each(function(){
					var score = $(this).data("score");
					var overlay = $('<div class="overlay" style="width:'+ (100 - parseInt(score)) +'%"></div>');
					$(this).prepend(overlay);
				})
			}
		},

		collapseBar: function(){
			if($(".collapse-bar").length){
				$(".collapse-bar .collapse-head").on("click", function(){
					if(!$(this).parent().hasClass("active")){
						$(this).parent().children(".collapse-body").slideDown(200);
						$(this).parent().addClass("active");
						$(this).find(".collapse-icon .icon").removeClass("ion-plus").addClass("ion-minus");
					}else{
						$(this).parent().children(".collapse-body").slideUp(200);
						$(this).parent().removeClass("active");
						$(this).find(".collapse-icon .icon").addClass("ion-plus").removeClass("ion-minus");
					}
				})
			}
		}
	}.init();
})(jQuery);