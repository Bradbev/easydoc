function tocClick(topUrl, contentUrl) {
    var $content = $('#content');
    $content.attr('src', contentUrl)
    window.location.replace(topUrl)
    return false
}

$(function() {
    // href clicking in the TOC
    var $content = $('#content');
    $('#toc a')
        .click(function(topUrl, contentUrl) {
            return false;
        });

    // searchbox input
    box = $('#searchBox')
    box.keyup(function(event) {
        // Number 13 is the "Enter" key on the keyboard
        if (event.keyCode === 13) {
            // Cancel the default action, if needed
            event.preventDefault();
            window.location.replace(getBase() + '/#/?search=' + box.val())
            window.location.reload()
        }
    });
    // populate search box from url if needed
    searchPrefix = "#/?search="
    if (window.location.hash.startsWith(searchPrefix)) {
        content = window.location.hash.slice(searchPrefix.length)
        box.val(content)
    }

    // iframe redirection on URL change
    $(document).ready(function() {
        var $content = $('#content');
        console.log(window.location.hash)
        if (window.location.hash != "") {
            noHash = window.location.hash.substring(1)
            $content.attr('src', noHash)
        }
    });
})