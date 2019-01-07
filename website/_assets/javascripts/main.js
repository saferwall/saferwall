//= require libraries/jquery
//= require libraries/materialize.min
//= require libraries/jquery.sticky
//= require libraries/highlight.min
//= require libraries/wow.min
//= require libraries/particles.min

(function($){
	return{
		init : function(){
			if($("select").length) $('select').material_select();
			if($(".tooltipped").length) $('.tooltipped').tooltip({delay: 50});
			$('pre code').each(function(i, block) {
				hljs.highlightBlock(block);
			});
			$("img").attr("width", "").attr("height", "");
			if($(".wow").length) new WOW().init();

			if($('.atf').length) Particles.init({selector: '.background'});
			
			this.advancedOptions();
			this.scrollDown();
			this.mobile();
			this.sticky();
			this.docSidebar();
			this.inputType();
			this.showFileName();
			// this.reports();
			this.stickySidebar()
			this.fileUpload()
		},

		inputType : function(){
			if($(".upload-source").length){
				$(".upload-source select").on("change", function(){
					if($(this).val() == "File"){
						$(".input-source .input-form").addClass("hidden");
						$(".input-source .input-form.file-source").removeClass("hidden");
					}else if ($(this).val() == "Url"){
						$(".input-source .input-form").addClass("hidden");
						$(".input-source .input-form.url-source").removeClass("hidden");
					}else{
						$(".input-source .input-form").addClass("hidden");
						$(".input-source .input-form.search-source").removeClass("hidden");
					}
				})
			}
		},

		showFileName: function(){
			if($(".upload-source").length){
				$(".input-form.file-source input").on("change", function(){
					var filename = $(this).val().split('\\').pop();
					$(this).parent().find("label").addClass("active").html(filename);
				})
			}
		},

		scrollDown : function(){
			if($(".scroll-down").length){
				$(".scroll-down").on("click", function(){
					$('html, body').animate({
						scrollTop: $(".benefits").offset().top
					}, 500);
				})
			}
		},

		advancedOptions : function(){
			if($(".advanced-options").length){
				this.showOptions();
				this.hideOptions();
			}
		},

		showOptions : function(){
			$(".advanced-button").on("click", function(){
				$(".advanced-options").slideDown(200);
			})
		},

		hideOptions : function(){
			$(".advanced-options").find(".hide-options").on("click", function(){
				$(this).parents(".advanced-options").slideUp(200);
			})
		},

		mobile : function(){
			this.fillMobileNav();
			this.toggleMobileNav();
		},

		fillMobileNav : function(){
			if($("nav.main-nav").length){
				var mobileNav = '<div class="mobile-nav">';
				mobileNav += $("nav.main-nav").html();
				mobileNav += '</div>';


				$("header.main-header").append(mobileNav);
			}
		},

		toggleMobileNav : function(){
			if($(".mobile-nav").length){
				$(".mobile-nav-icon").on("click", function(){
					$(".mobile-nav").slideToggle(200);
				})
			}
		},

		sticky : function(){
			if($(".vertical-nav").length){
				if($(window).width() > 768){
					$(".vertical-nav").sticky({
						topSpacing: 20,
						bottomSpacing : $("footer.main-footer").height() + 100
					});
				}
			}
		},

		docSidebar : function(){
			if($(".doc-sidebar").length){
				this.initDoc();
				this.stickyDocSidebar();
				this.docToggleMenu();
				this.moveToSection();
			}
		},

		docToggleMenu : function(){
			$(".doc-sidebar .nav-section .section-title").on("click", function(){
				if(!$(this).hasClass("active")){
					$(this).addClass("active");
					$(this).parents(".nav-section").find("ul").slideDown(100);
				}else{
					$(this).removeClass("active");
					$(this).parents(".nav-section").find("ul").slideUp(100);
				}
			})
		},

		stickyDocSidebar : function(){
			if($(window).width() > 768){
				$(".doc-sidebar").sticky({
					topSpacing: 20,
					bottomSpacing: $("footer.main-footer").height() + 100
				})
			}
		},

		moveToSection: function(){
			$(".doc-sidebar .nav-section ul li a").on("click", function(e){
				e.preventDefault();

				$(this).parents(".doc-sidebar").find("ul li a.active").removeClass("active");
				$(this).addClass("active");

				var dest = $(this).attr("href").replace("#", "");
				$('html, body').animate({
					scrollTop: $("#" + dest).offset().top - 20
				}, 1000);
			})
		},

		initDoc : function(){
			$(".doc-sidebar .nav-section.active ul").slideDown();
		},

		stickySidebar(){
			if($('.blog-sidebar').length){
				if($(window).width() > 768){
					$(".blog-sidebar").sticky({
						topSpacing: 20,
						bottomSpacing: $("footer.main-footer").height() + 100
					})
				}
			}
		},

		fileUpload(){
			if($('.file-upload').length){
				$(".file-upload input[type='file']").on("change", function(e){
					let file = e.target.files[0],
						filename = file.name
					$(this).parent().find(".filename").html(filename)
					if(file.size > 64000000) { 
						this.notificationError = "over sized."
						this.notifActive = true
						return
					}

					var reader = new FileReader();
					var hashcode = ''
					reader.onload = (function(theFile) {
						return function(e) {
							hashcode = sha256(e.target.result)
							let url = Global.apiUrl + hashcode + '?api-key=' + Global.apiKey
							$.ajax({
								method: 'get',
								url: url,
								data: {},
								success: function(data){
									console.log(data)
								}
							})
						};
					})(file);
					reader.readAsDataURL(file);
				})
			}
		}
		
	}.init();
})(jQuery);