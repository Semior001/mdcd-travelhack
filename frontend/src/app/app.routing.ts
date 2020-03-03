import {Routes} from '@angular/router';

import {FullComponent} from './layouts/full/full.component';
import {MainComponent} from './main/main.component';
import {UsersComponent} from './users/users.component';
import {LoginComponent} from './auth/login/login.component';
import {ResetPasswordComponent} from './auth/reset-password/reset-password.component';
import {AuthComponent} from './auth/auth.component';
import {UserComponent} from './user/user.component';
import {NotFoundComponent} from './not-found/not-found.component';
import {BackgroundsComponent} from './backgrounds/backgrounds.component';
import {AuthGuard} from './auth/auth.guard';

// export const AppRoutes: Routes = [
//   {
//     path: 'auth', component: AuthComponent, children: [
//       {path: 'login', component: LoginComponent},
//       {path: 'reset-password', component: ResetPasswordComponent},
//       {path: '**', redirectTo: '/not-found', pathMatch: 'full'}
//     ]
//   },
//   {path: 'not-found', component: NotFoundComponent},
//   {
//     path: '', component: FullComponent,
//     children: [
//       {path: 'admin/users/:id', component: UserComponent},
//       {path: 'admin/users', component: UsersComponent},
//       {path: 'admin/backgrounds', component: BackgroundsComponent},
//       {path: 'main', component: MainComponent},
//       {path: '', redirectTo: 'main', pathMatch: 'full'},
//       {path: '**', redirectTo: 'not-found', pathMatch: 'full'}
//     ]
//   },
//   {path: '**', redirectTo: '/not-found', pathMatch: 'full'}
// ];

export const AppRoutes: Routes = [
  {
    path: 'auth', component: AuthComponent, children: [
      {path: 'login', component: LoginComponent},
      {path: 'reset-password', component: ResetPasswordComponent},
      {path: '**', redirectTo: '/not-found', pathMatch: 'full'}
    ]
  },
  {path: 'not-found', component: NotFoundComponent},
  {
    path: '', component: FullComponent, canActivate: [AuthGuard],
    children: [
      {path: 'admin/users/:id', component: UserComponent, canActivate: [AuthGuard]},
      {path: 'admin/users', component: UsersComponent, canActivate: [AuthGuard]},
      {path: 'admin/backgrounds', component: BackgroundsComponent, canActivate: [AuthGuard]},
      {path: 'main', component: MainComponent, canActivate: [AuthGuard]},
      {path: '', redirectTo: 'main', pathMatch: 'full', canActivate: [AuthGuard]},
      {path: '**', redirectTo: 'not-found', pathMatch: 'full'}
    ]
  },
  {path: '**', redirectTo: '/not-found', pathMatch: 'full'}
];
