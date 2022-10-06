@extends('layouts.main')

@section('content')
<script src="//ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>

<script>
$(document).ready(function() {

  $('#loginform').submit(function() {

      $.ajax({
          type: "GET",
          url: 'api/login',
          data: {
              password: $("#password").val()
          },
          success: function(data)
          {
              if (data === 'Login Successful') {
                  window.location.replace('/index');
              }
              else {
                (document.getElementById('alert')).style.visibility = 'visible';
                document.getElementById('alert').innerHTML = 'Invalid Login';

              }
          }
      });
      return false;
  });
});
</script>
<body class="text-center">

  <main class="form-signin">

  <div class="mask d-flex align-items-center h-100 gradient-custom-3">
    <div class="container h-100">
      <div class="row d-flex justify-content-center align-items-center h-100">
        <div class="col-12 col-md-9 col-lg-7 col-xl-6">
          <div class="card" style="border-radius: 15px;">
            <div class="card-body p-5">
              <h2 class="text-uppercase text-center mb-5">Login</h2>
              <p>Please enter a password in order to get the flag</p>

                    <div id="alert" style="visibility: hidden;" class="alert alert-info"></div>

              <form id="loginform" name="loginform">
                <div class="form-outline mb-4">
                  <input type="password" name="password" id="password" class="form-control form-control-lg" />
                </div>

                <div class="d-flex justify-content-center">
                  <button id="login" type="submit" class="btn btn-outline-success">Login</button>
                </div>

              </form>

            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

    </main>


</body>
