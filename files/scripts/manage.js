$(document).ready(function(){
	$("#submit_addteam").click(function(){
		$("#form_addteam").submit();
	});

	$(".delete-team").click(function(e){
		e.preventDefault();
		var txt;
		var conf = confirm("Are you sure?");
		if (conf) {
		    $("#form-delete-team-" + $(this).attr("data-team-id")).submit();
		}
	});

	$(".delete-member").click(function(e){
		e.preventDefault();
		var txt;
		var conf = confirm("Are you sure?");
		if (conf) {
		    $("#form-delete-member-" + $(this).attr("data-member-id")).submit();
		}
	});
})