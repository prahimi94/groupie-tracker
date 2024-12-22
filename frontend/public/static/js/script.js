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

const first_album_date_start = document.getElementById('first_album_date_start');
const first_album_date_start_value = document.getElementById('first_album_date_start_value');
const first_album_date_end = document.getElementById('first_album_date_end');
const first_album_date_end_value = document.getElementById('first_album_date_end_value');

const select2Data = JSON.parse(allUniqueLocations).sort().map(item => ({
  id: item, // The value for the <option>
  text: item.replace(/_/g, ' ').replace('-', ', ') // Display text with formatted replacements
}));

const $select = $('#concerts_locations');
$select.select2({
  data: select2Data,
  placeholder: 'Select locations',
  allowClear: true // Allows user to clear the selection
});

$("#concerts_locations").val('').change();

function filter_result() {

  creation_date_end_value.textContent = creation_date_end.value;
  creation_date_start_value.textContent = creation_date_start.value;
  first_album_date_start_value.textContent = first_album_date_start.value;
  first_album_date_end_value.textContent = first_album_date_end.value;

  const selectedLocations = $('#concerts_locations').val(); // Get selected values as an array
  // Get all currently checked values
  const checkboxes = document.querySelectorAll('input[name="members[]"]');
  const selectedMemberCounts = Array.from(checkboxes)
    .filter(checkbox => checkbox.checked)
    .map(checkbox => parseInt(checkbox.value, 10));

  $.each(JSON.parse(allArtists), function( index, value ) {
    const dateString = value.firstAlbum
    const parts = dateString.split("-");
    const year = parts[2]; 
    const membersCount = value['members'].length;
    var showArtistForLocationFilter = true

    if (selectedLocations && selectedLocations.length > 0) {
      showArtistForLocationFilter = false
      $.each(value['LocationsData'], function(index2, location) {
        location.replace(/_/g, ' ').replace('-', ', ')
        if (location == selectedLocations){
          showArtistForLocationFilter = true;
        }
      }) 
    }
    
    if(value.creationDate >= creation_date_start.value && value.creationDate <= creation_date_end.value
      && year >= first_album_date_start.value && year <= first_album_date_end.value
      && showArtistForLocationFilter
      && selectedMemberCounts.includes(membersCount)
    ) {
      $('#artist_' + value.id).show()
    } else {
      $('#artist_' + value.id).hide()
    }
  });
}

function resetForm(){

  $("#concerts_locations").val('').change();

  $('input[name="members[]"]').prop('checked', true);

  $('#creation_date_start').val(1950).change();
  $('#creation_date_end').val(2020).change();
  $('#first_album_date_start').val(1950).change();
  $('#first_album_date_end').val(2020).change();

  $('[id^="artist"]').show();
}