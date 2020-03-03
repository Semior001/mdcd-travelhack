import {
  CanActivate,
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
  Router
} from '@angular/router';
import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';

import {AuthService} from './auth.service';

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(
    private authService: AuthService,
    private router: Router) {
  }

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): boolean | Observable<boolean> | Promise<boolean> {
    console.log('Guard called');
    // const isAuth = this.authService.isAuthenticated();
    // if (!isAuth) {
    //   this.router.navigate(['/auth/login']);
    // }
    // } else {
    //   const role = this.authService.getRole();
    //   if (role === 'admin') {
    //     return true;
    //   } else if (role === 'staff') {
    //     return !state.url.startsWith('/admin');
    //   }
    // }
    // return this.authService.isAuthenticated();
    // return this.authService.isAuthenticated();
    return true;
  }
}
