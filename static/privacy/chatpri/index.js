$(function() {

   let currKey = localStorage.getItem('oldKey')
   if(currKey){
      window.location.href = "/chatIndex?kefu_id=caonima888&ent_id=1&visitor_name="+currKey
   }else{
      $("#privePass").val("")

      $(".input input").focus(function() {

         $(this).parent(".input").each(function() {
            $("label", this).css({
               "line-height": "18px",
               "font-size": "18px",
               "font-weight": "100",
               "top": "0px"
            })
            $(".spin", this).css({
               "width": "100%"
            })
         });
      }).blur(function() {
         $(".spin").css({
            "width": "0px"
         })
         if ($(this).val() == "") {
            $(this).parent(".input").each(function() {
               $("label", this).css({
                  "line-height": "60px",
                  "font-size": "24px",
                  "font-weight": "300",
                  "top": "10px"
               })
            });

         }
      });

      $(".button").click(function(e) {


         if(!$("#privePass").val()){

            // alert('需填写秘钥')
            Qmsg.error("需填写秘钥!",{});
            return false

         }


         // var pX = e.pageX,
         //     pY = e.pageY,
         //     oX = parseInt($(this).offset().left),
         //     oY = parseInt($(this).offset().top);
         //
         // $(this).append('<span class="click-efect x-' + oX + ' y-' + oY + '" style="margin-left:' + (pX - oX) + 'px;margin-top:' + (pY - oY) + 'px;"></span>')
         // $('.x-' + oX + '.y-' + oY + '').animate({
         //    "width": "500px",
         //    "height": "500px",
         //    "top": "-250px",
         //    "left": "-250px",
         //
         // }, 600);
         // $("button", this).addClass('active');
         let tempLocatStr = $("#privePass").val()
         $("#privePass").val("")

         let requestdata = {}
         requestdata.key = tempLocatStr
         $.ajax({
            url: "/JyKey",
            data: requestdata,
            type: "POST",
            dataType: "json",
            success: function (resultJson) {

               if (200 === resultJson.code) {

                  localStorage.setItem("oldKey",tempLocatStr);
                  window.location.href = "/chatIndex?kefu_id=caonima888&ent_id=1&visitor_name="+tempLocatStr

               }else{
                  // alert('请联系客服提供密钥!')
                  var configs = {};
                  // configs 为配置参数，可省略
                  Qmsg.error("请联系客服提供密钥!",configs);
               }


            },


         });



      })

      $(".alt-2").click(function() {
         if (!$(this).hasClass('material-button')) {
            $(".shape").css({
               "width": "100%",
               "height": "100%",
               "transform": "rotate(0deg)"
            })

            setTimeout(function() {
               $(".overbox").css({
                  "overflow": "initial"
               })
            }, 600)

            $(this).animate({
               "width": "140px",
               "height": "140px"
            }, 500, function() {
               $(".box").removeClass("back");

               $(this).removeClass('active')
            });

            $(".overbox .title").fadeOut(300);
            $(".overbox .input").fadeOut(300);
            $(".overbox .button").fadeOut(300);

            $(".alt-2").addClass('material-buton');
         }

      })

      $(".material-button").click(function() {

         if ($(this).hasClass('material-button')) {
            setTimeout(function() {
               $(".overbox").css({
                  "overflow": "hidden"
               })
               $(".box").addClass("back");
            }, 200)
            $(this).addClass('active').animate({
               "width": "700px",
               "height": "700px"
            });

            setTimeout(function() {
               $(".shape").css({
                  "width": "50%",
                  "height": "50%",
                  "transform": "rotate(45deg)"
               })

               $(".overbox .title").fadeIn(300);
               $(".overbox .input").fadeIn(300);
               $(".overbox .button").fadeIn(300);
            }, 700)

            $(this).removeClass('material-button');

         }

         if ($(".alt-2").hasClass('material-buton')) {
            $(".alt-2").removeClass('material-buton');
            $(".alt-2").addClass('material-button');
         }

      });

   }


});
