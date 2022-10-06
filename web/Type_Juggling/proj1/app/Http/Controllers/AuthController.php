<?php

namespace App\Http\Controllers;

use App\Models\User;
use Illuminate\Http\Request;
use App\Http\Requests\RegisterRequest;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Hash;



class AuthController extends Controller
{
    /**
     * Display login page.
     *
     * @return \Illuminate\Http\Response
     */
    public function show_login()
    {
        return view('auth.login');
    }



    /**
     * Handle account login
     *
     */
    public function customLogin(Request $request)
    {
        $request->validate([
            'password' => 'required',
        ]);

        if ($request->get('password') == "Bsid3s_k3ny4_hihi!") {
            session(['loggedin' => True]);
            return "Login Successful";
        }

        return "Invalid Password";
    }

}