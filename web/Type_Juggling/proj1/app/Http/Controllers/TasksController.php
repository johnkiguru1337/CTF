<?php

namespace App\Http\Controllers;

use Illuminate\Support\Facades\Auth;
use Illuminate\Http\Request;
use App\Models\Task;
use DB;



class TasksController extends Controller
{
    //

    public function index(request $request) {
        if ($request->session()->has('loggedin')) {
            return view('/index');
            die();
        } else {
            return redirect('/login');
            die();
        }
    }
}
