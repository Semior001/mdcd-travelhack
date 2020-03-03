import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaderResponse, HttpResponse} from '@angular/common/http';
import {Router} from '@angular/router';

import {environment} from '../../environments/environment';
import {MatSnackBar} from '@angular/material';
import {CookieService} from 'ngx-cookie-service';

import { Cookie } from 'ng2-cookies/ng2-cookies';

const BACKEND_URL = environment.apiUrl;

@Injectable({providedIn: 'root'})
export class AuthService {
  isLogged = false;

  constructor(
    private http: HttpClient,
    private router: Router,
    private snackBar: MatSnackBar,
    private cookiesService: CookieService
  ) {
  }

  isAuthenticated(): boolean {
    return this.isLogged;
  }

  createUser(email: string, password: string) {
    // this.http.post(BACKEND_URL + "/signup", {email, password}).subscribe(
    //   () => {
    //     this.router.navigate(["/"]);
    //   },
    //   error => {
    //     this.authStatusListener.next(false);
    //   }
    // );
  }

  autoLogin() {
    if (this.cookiesService.get('JWT') !== null) {
      this.isLogged = true;
    }
  }

  login(email: string, password: string) {
    this.http
      .get(
        BACKEND_URL + `/auth/local/login?user=${email}&passwd=${password}&aud=mdcd_api_middleware&session=1`,
        {observe: 'response', withCredentials: true})
      .subscribe(
        (response) => {
          this.router.navigate(['/']);
          this.isLogged = true;
          console.log('Logged in');
        },
        (error) => {
          console.log(error);
          if (error === 'failed to check user credentials') {
            this.snackBar.open('Неверный логин или пароль');
          } else {
            this.snackBar.open('Неизвестная ошибка');
          }
        }
      );
  }

  logout() {
    this.http.get(
      BACKEND_URL + `/auth/local/logout`
      ,
      {withCredentials: true}).subscribe(
      () => {
        this.isLogged = false;
        sessionStorage.clear();

        sessionStorage.setItem('JWT', '');
        Cookie.deleteAll();

        this.router.navigate(['/auth/login']);
        console.log('Logged out');
      },
      () => {
        this.snackBar.open('Неизвестная ошибка');
      }
    );
    // this.cookiesService.delete('JWT');
    // this.router.navigate(['/auth/login']);
    // console.log('Logged out');
  }
}
