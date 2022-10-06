<?php

use App\Http\Controllers\TasksController;
use App\Http\Controllers\AuthController;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/
Route::get('/' ,[AuthController::class, 'show_login'])->name('login');
Route::get('/index', [TasksController::class, 'index']);
Route::get('/login', [AuthController::class, 'show_login'])->name('login');