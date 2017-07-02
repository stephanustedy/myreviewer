function init_form() {
	
	var selectedTeam = $("#team").val();

	//cleanup
	$('#qa').find('option').remove();
	$('#se').find('option').remove();
	$("#pr").val('');

	for(var i=0 ;i < team[selectedTeam].Member.length ;i++) {
		var member = team[selectedTeam].Member[i];
		if(member.Role == "QA") {
			$("#qa").append($('<option>', {
			    value: member.MemberId,
			    text: member.Name
			}));
		}
		if(member.Role == "SE") {

			$("#se").append($('<option>', {
			    value: member.MemberId,
			    text: member.Name
			}));
		}
	}
}

$(document).ready(function(){
	init_form();
	$("#team").change(function(){
		init_form();
	});

	$("#submit-review").click(function(e){
		e.preventDefault();
		if(!$("#team").val() || !$("#qa").val() || !$("#se").val() || !$("#pr").val()) {
			alert("Fill all forms!");
		} else {
			$("#review-form").submit();
		}
	});
});