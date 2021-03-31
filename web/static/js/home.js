function postPaste(content, ttl) {
    jQuery.post({
        url: "/api/v1/paste",
        contentType: 'application/json',
        processData: false,
        data: JSON.stringify({
            content: content,
            ttl: ttl
        }),
        error: function (jqXHR, textStatus, errorThrown) {
            var response = JSON.parse(jqXHR.responseText);
            for (i in response['errors']) {
                jQuery('<div class="alert alert-dismissible alert-danger"><button type="button" class="close" data-dismiss="alert" onclick="javascript:jQuery(this).parent().remove();">&times;</button><strong id="error-msg">' + response['errors'][i] + '</strong></div>').insertBefore('div.form-wrapper');
            }
        },
        success: function (data, textStatus, jqXHR) {
            jQuery('div.form-wrapper').hide();
            jQuery('div.result-wrapper p.card-text input').attr('value', location.href + data['recovery_key']);
            jQuery('div.result-wrapper').show();
        }
    });
}

function copyToClipboard(el) {
    console.log(el);
    /* Select the text field */
    el.focus();
    el.select();
    el.setSelectionRange(0, 99999); /* For mobile devices */

    /* Copy the text inside the text field */
    document.execCommand("copy");
}

function downloadFile(fileUrl) {
    var link = document.createElement('a')
    link.href = fileUrl;
    link.setAttribute('download', true);
    link.click();
    link.remove();
}

jQuery(document).ready(function () {
    jQuery('#btn-paste-view').click(function(){
        var recoveryKey = jQuery('#recovery-key').val();
        if (recoveryKey != "") {
            var re = /^data:(.*);base64,/;
            jQuery.getJSON({
                url: "/api/v1/paste/"+recoveryKey,
                error: function (jqXHR, textStatus, errorThrown) {
                    var response = JSON.parse(jqXHR.responseText);
                    for (i in response['errors']) {
                        jQuery('<div class="alert alert-dismissible alert-danger"><button type="button" class="close" data-dismiss="alert" onclick="javascript:jQuery(this).parent().remove();">&times;</button><strong id="error-msg">' + response['errors'][i] + '</strong></div>').insertBefore('#paste-warning');
                    }
                },
                success: function (data, textStatus, jqXHR) {
                    var res = re.exec(data['content'])
                    if (res) {
                        /*var contentType = res[1];
                        var content = data['content'].replace(re, '');
                        var d = [];
                        d.push(atob(content));
                        var properties = {type: contentType};
                        var file = new Blob(d, properties);
                        var url = URL.createObjectURL(file);*/

                        downloadFile(data['content']);
                        document.location.href="/"
                    }
                    else {
                        jQuery('#paste-warning').hide();
                        jQuery('#paste-file-upload').hide();
                        jQuery('div.paste-option').hide();
                        jQuery('#paste-content-copy').show();
                        jQuery('#content').prop('readonly', true).text(data['content']);
                        jQuery('div.form-wrapper').show();
                    }
                }
            });
        }

        return false;
    });

    jQuery('#paste-create-form #btn-paste-submit').click(function (e) {
        e.preventDefault();

        postPaste(jQuery('#content').val(), jQuery('#expiration').val());

        return false;
    });

    jQuery('#file-upload').change(function(){
        const file = this.files[0];
        const reader = new FileReader();

        jQuery(reader).on("load", function(){
            //jQuery('#content').text(reader.result);

            postPaste(reader.result, jQuery('#expiration').val());
        });

        if (file) {
            reader.readAsDataURL(file);
        }
    });

    jQuery('#link-paste-url-copy').click(function (e) {
        e.preventDefault();
        var copyText = document.getElementById('paste-url');

        copyToClipboard(copyText);

        return false;
    });

    jQuery('#btn-paste-content-copy').click(function (e) {
        e.preventDefault();
        var copyText = document.getElementById('content');

        copyToClipboard(copyText);

        return false;
    });
})