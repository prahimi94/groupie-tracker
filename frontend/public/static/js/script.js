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

const creation_date_start = document.getElementById('creation_date_start');
const creation_date_start_value = document.getElementById('creation_date_start_value');
const creation_date_end = document.getElementById('creation_date_end');
const creation_date_end_value = document.getElementById('creation_date_end_value');

creation_date_end.addEventListener('input', () => {
  creation_date_end_value.textContent = creation_date_end.value;
  check_dates();
});

creation_date_start.addEventListener('input', () => {
  creation_date_start_value.textContent = creation_date_start.value;
  check_dates();
});

const first_album_date_start = document.getElementById('first_album_date_start');
const first_album_date_start_value = document.getElementById('first_album_date_start_value');
const first_album_date_end = document.getElementById('first_album_date_end');
const first_album_date_end_value = document.getElementById('first_album_date_end_value');
first_album_date_start.addEventListener('input', () => {
  first_album_date_start_value.textContent = first_album_date_start.value;
  check_dates();
});

first_album_date_end.addEventListener('input', () => {
  first_album_date_end_value.textContent = first_album_date_end.value;
  check_dates();
});

function check_dates() {
  $.each(JSON.parse(allArtists), function( index, value ) {
    const dateString = value.firstAlbum
    const parts = dateString.split("-");
    const year = parts[2]; 
    if(value.creationDate >= creation_date_start.value && value.creationDate <= creation_date_end.value
      && year >= first_album_date_start.value && year <= first_album_date_end.value
    ) {
      $('#artist_' + value.id).show()
    } else {
      $('#artist_' + value.id).hide()
    }
  });
}

const checkboxes = document.querySelectorAll('input[name="members[]"]');
checkboxes.forEach(checkbox => {
  checkbox.addEventListener('change', () => {
    const isChecked = checkbox.checked;
    const value = checkbox.value;

    // Get all currently checked values
    const selectedValues = Array.from(checkboxes)
      .filter(checkbox => checkbox.checked)
      .map(checkbox => parseInt(checkbox.value, 10));

    $.each(JSON.parse(allArtists), function( index, value ) {
        const membersCount = value['members'].length;

        if(!selectedValues.includes(membersCount)){
          $('#artist_' + value.id).hide()
        } else {
          $('#artist_' + value.id).show()
        }
    });
  });
});

const select2Data = JSON.parse(allUniqueLocations).map(item => ({
  id: item, // The value for the <option>
  text: item.replace(/_/g, ' ').replace('-', ', ') // Display text with formatted replacements
}));

$('.js-example-basic-multiple').select2({
  data: select2Data
});