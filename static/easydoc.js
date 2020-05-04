// Plugged into <a> elements from the table of contents menu
function tocClick(topUrl, contentUrl) {
    var $content = $('#content');
    $content.attr('src', contentUrl)
    window.location.replace(topUrl)
    return false
}

function setCssOpenClass($el, isOpen) {
    if (isOpen) {
        $el.removeClass('closed');
        $el.addClass('open');
    } else {
        $el.removeClass('open');
        $el.addClass('closed');
    }
}

function toggleCssOpenClass($el) {
    if ($el.hasClass('open')) {
        setCssOpenClass($el, false);
    } else {
        setCssOpenClass($el, true);
    }
}

$(function() {
    var $content = $('#content');
    var $toc = $('#toc');
    var $tocDrawerToggle = $('#toc-drawer-toggle')
    var $searchBox = $('#search-field');
    var $searchButton = $('#search-field-button');

    addMediaSpecificBehavior();

    // href clicking in the TOC
    $toc.find('a').click(function(topUrl, contentUrl) {
        return false;
    })
    $tocDrawerToggle.click(function(event) {
        toggleCssOpenClass($toc);
    });
    // searchbox input
    var executeSearch = () => {
        window.location.replace(getBase() + '/#/?search=' + $searchBox.val());
        window.location.reload();
    };

    $searchBox.keyup(function(event) {
        // Number 13 is the "Enter" key on the keyboard
        if (event.keyCode === 13) {
            // Cancel the default action, if needed
            event.preventDefault();
            executeSearch();
        }
    });
    $searchBox.focus(function(event) {
        // Prevents our button from animating on page load if search text is already present.
        if (!$searchButton.hasClass('animate')) {
            $searchButton.addClass('animate');
        }
    });
    $searchButton.click(function(event) {
        executeSearch();
    });
    // populate search box from url if needed
    var searchPrefix = "#/?search="
    if (window.location.hash.startsWith(searchPrefix)) {
        content = window.location.hash.slice(searchPrefix.length)
        $searchBox.val(content)
    }



    // iframe redirection on URL change
    $(document).ready(function() {
        var $content = $('#content');
        console.log(window.location.hash)
        if (window.location.hash != "") {
            noHash = window.location.hash.substring(1)
            $content.attr('src', getBase() + noHash)
        }
    });

    function mediaLargeBehavior() {
        setCssOpenClass($toc, true);
        $toc.addClass('push');
        var $tocDrawerDragbar = $('#toc-drawer-dragbar');
        var hasTouchEvents = 'ontouchstart' in document;
        function dragStart() {
            $toc.css('transition', 'unset');
            // Mousemove gets intercepted by the iframe
            // This covers the iframe so that while we are dragging the mouse events
            // keep flowing.
            $tocDrawerDragbar.width('100vw');
        }
        function dragEnd() {
            $toc.css('transition', '');
            $tocDrawerDragbar.width('');
        }
        function dragResize(event) {
            $toc.width(event.pageX);
        }
        if(hasTouchEvents) {
            $tocDrawerDragbar.on('touchstart', function() {
                dragStart();
                $(document).on('touchmove', dragResize);
            });
            $(document).on('touchend', function() {
                dragEnd();
                $(document).off('touchmove', dragResize);
            });
        } else {
            $tocDrawerDragbar.on('mousedown', function() {
                dragStart();
                $(document).on('mousemove', dragResize);
            });
            $(document).on('mouseup', function() {
                dragEnd();
                $(document).off('mousemove', dragResize);
            });
        }
    }

    function mediaSmallBehavior() {
        setCssOpenClass($toc, false);
        $toc.addClass('float');
        $toc.find('a').click(function() {
            toggleCssOpenClass($toc);
        });

        // Click outside closes the menu
        $content.load(function() {
            var $iframeDocument = $($content[0].contentDocument || $content[0].contentWindow.document);
            $iframeDocument.click(function(event) {
                setCssOpenClass($toc, false);
            });
        });
    }

    function addMediaSpecificBehavior() {
        if (window.matchMedia("(min-width: 900px)").matches) {
            mediaLargeBehavior();
        } else {
            mediaSmallBehavior();
        }
    }
})
