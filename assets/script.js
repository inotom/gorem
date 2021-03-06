(function($) {
  'use strict';

  var DEFAULT_FONT_SIZE = 14;

  var $btn = $('#btn-create'),
      $thumb = $('#thumb'),
      $fontSize = $('#fs'),
      $fontSizeOut = $('#fs-out'),
      $width = $('#w'),
      $height = $('#h'),
      $text = $('#s'),
      $hasProp = $('#p');

  $(function() {

    $fontSize.val(DEFAULT_FONT_SIZE);
    $fontSize.on('input', function() {
      var fs = $(this).val();
      $fontSizeOut.text(fs);
    });

    $btn.on('click', function() {
      $thumb.hide();
      var $img = $('<img>');
      var imgURI = '/lorem?w=' + $width.val() + '&h=' + $height.val() + '&fs=' + $fontSize.val();
      var text = $text.val();
      if (text.length > 0) {
      imgURI += '&s=' + text.replace(/\s/g, '+')
                            .replace(/&/g, '＆')
                            .replace(/=/g, '＝')
                            .replace(/;/g, '；')
                            .replace(/%/g, '％');
      }
      var hasProp = $hasProp.prop('checked');
      if (hasProp) {
        imgURI += '&p=1';
      }
      var imgType = $('input[name="t"]:checked').val();
      if (imgType) {
        imgURI += '&t=' + imgType;
      }
      $img.attr('src', imgURI);
      $thumb
        .empty()
        .append($img)
        .fadeIn();
    });
  });
})(jQuery);
