import {Injectable} from '@angular/core';

export interface Menu {
  state: string;
  name: string;
  type: string;
  icon: string;
}

const MENUITEMS = [
  // {state: '/auth/login', name: 'Логин', type: 'link', icon: 'av_timer'},
  // {state: '/auth/reset-password', name: 'Сброс пароля', type: 'link', icon: 'av_timer'},
  {state: '/main', name: 'Главная', type: 'link', icon: 'av_timer'},
  {state: '/admin/backgrounds', name: 'Добавление фонов', type: 'link', icon: 'av_timer'},
  // {state: '/admin/users', name: 'Пользователи', type: 'link', icon: 'av_timer'},
];

@Injectable()
export class MenuItems {
  getMenuitem(): Menu[] {
    return MENUITEMS;
  }
}
