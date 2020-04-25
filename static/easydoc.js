$(function() {
    var $content = $('#content');
    $('#toc a')
        .click(function() {
            $content.attr('src', $(this).attr('href'));
            window.location.hash = $(this).attr('href');
            return false;
        });
});

$(document).ready(function() {
    var $content = $('#content');
    console.log(window.location.hash)
    if (window.location.hash != "") {
        noHash = window.location.hash.substring(1)
        $content.attr('src', noHash)
    }
});