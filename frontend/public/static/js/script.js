var Swipes = new Swiper('.swiper-container', {
    spaceBetween: 0,
    centeredSlides: false,
    speed: 5000,
    autoplay: {
    delay: 1,
    },
    loop: true,
    loopedSlides: 4,
    slidesPerView:'auto',
    allowTouchMove: false,
    disableOnInteraction: true
});
// $(document).ready(function(){
// $("#find-artist").on("click", function() {
//     $('.artist-name').each(function(i, obj) {
//         var value = $("#search-text").val().toLowerCase();
//         var element_p_id = $(this).attr('id');
//         var element_a_id = element_p_id.replace("-p-", "-a-");
//         if ($(this).text().toLowerCase().indexOf(value) > -1) {
//         console.log($(this).text())
//         console.log(element_a_id)
//         $("#"+element_a_id).show();
//         } else {
//         $("#"+element_a_id).hide();
//         }
//     });
//     });
// });

const creation_date = document.getElementById('creation_date');
const creation_date_value = document.getElementById('creation_date_value');
if (creation_date) {
    creation_date.addEventListener('input', () => {
    creation_date_value.textContent = creation_date.value;

    $.each(JSON.parse(allArtists), function( index, value ) {
        if(value.creationDate > creation_date.value) {
        $('#artist_' + value.id).hide()
        } else {
        $('#artist_' + value.id).show()
        }
    });
    
    });
}


const first_album_date = document.getElementById('first_album_date');
const first_album_date_value = document.getElementById('first_album_date_value');
if(first_album_date){
    first_album_date.addEventListener('input', () => {
    first_album_date_value.textContent = first_album_date.value;

    $.each(JSON.parse(allArtists), function( index, value ) {
        const dateString = value.firstAlbum
        const parts = dateString.split("-");
        const year = parts[2];  
        if(year > first_album_date.value) {
        $('#artist_' + value.id).hide()
        } else {
        $('#artist_' + value.id).show()
        }
    });

    });
}
    

var yourNavigation = $(".nav-filter");
    stickyDiv = "sticky-filter";
    yourHeader = $('.hero-title').outerHeight() + $('nav.navbar').outerHeight() + 150;

$(window).scroll(function() {
  if( $(this).scrollTop() > yourHeader ) {
    yourNavigation.addClass(stickyDiv);
  } else {
    yourNavigation.removeClass(stickyDiv);
  }
});

