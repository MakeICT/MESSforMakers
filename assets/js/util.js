document.addEventListener('DOMContentLoaded', function() {
  var elems = document.querySelectorAll('.datepicker');
  var options = {"disableWeekends":true};
  var instances = M.Datepicker.init(elems, options);
});