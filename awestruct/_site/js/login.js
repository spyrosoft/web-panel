document.getElementById('password').focus();

success_response_callbacks['login'] = function(post_data, results) {
	window.location = '/';
};
