//= require libraries/jquery

(function($){
	return{
		init : function(){
			this.showPassword();
			this.fieldFocus();
		},

		showPassword : function(){
			$(".preview").on("click", function(){
				$(this).toggleClass("shown");

				if(!$(this).hasClass("shown"))
					$(this).parents(".sign-field").find("input").attr("type", "text");
				else
					$(this).parents(".sign-field").find("input").attr("type", "password");
			})
		},

		fieldFocus : function(){
			$(".sign-field input").on("focus", function(){
				$(this).parents(".sign-field").addClass("active");
			})
			$(".sign-field input").on("blur", function(){
				$(this).parents(".sign-field").removeClass("active");
			})
		}
	}.init();
})(jQuery);