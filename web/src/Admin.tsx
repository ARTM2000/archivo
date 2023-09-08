import {
  Admin,
  Resource,
  useNotificationContext,
  useTranslate,
} from 'react-admin';
import { DataProvider } from './data-provider';
import { AuthProvider, PERMISSIONS } from './auth-provider';
import { Dashboard } from './components/home/dashboard';
import { LoginWrapper } from './components/auth/login-wrapper';
import { SourceServerList } from './components/sourceservers/list';
import { SourceServerCreate } from './components/sourceservers/create';
import StorageSharpIcon from '@mui/icons-material/StorageSharp';
import PeopleAltSharpIcon from '@mui/icons-material/PeopleAltSharp';
import { FilesList } from './components/files/list';
import { Route } from 'react-router-dom';
import { FileSnapshotsShow } from './components/files/show';
import { UserList } from './components/users/list';
import { RegisterUser } from './components/users/register-user';
import { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { ShowUserActivities } from './components/users/show-activites';

export const AdminPanel = () => {
  const [perm, setPerm] = useState<PERMISSIONS>(PERMISSIONS.USER);
  useEffect(() => {
    AuthProvider.getPermissions('').then((perm: PERMISSIONS) => {
      setPerm(perm);
    });
    const permissionCheckInterval = setInterval(
      () =>
        AuthProvider.getPermissions('').then((perm: PERMISSIONS) => {
          setPerm(perm);
        }),
      3000,
    );
    return () => clearInterval(permissionCheckInterval);
  }, []);

  return (
    <Admin
      dataProvider={DataProvider as any}
      authProvider={AuthProvider}
      dashboard={Dashboard}
      loginPage={LoginWrapper}
      notification={CustomAdminPanelNotification}
      requireAuth
    >
      <Resource
        name="servers"
        list={SourceServerList}
        create={SourceServerCreate}
        hasEdit={false}
        icon={StorageSharpIcon}
      >
        <Route path=":serverId/:serverName/files" element={<FilesList />} />
        <Route
          path=":serverId/:serverName/files/:filename"
          element={<FileSnapshotsShow />}
        />
      </Resource>
      {perm === PERMISSIONS.ADMIN && (
        <Resource
          name="users"
          list={UserList}
          create={RegisterUser}
          icon={PeopleAltSharpIcon}
        >
          <Route path=":userId/:username/activities" element={<ShowUserActivities />} />
        </Resource>
      )}
    </Admin>
  );
};

const CustomAdminPanelNotification = (...props: any[]) => {
  const { notifications, takeNotification } = useNotificationContext();
  const [messageInfo, setMessageInfo] = useState<any>(null);
  const [open, setOpen] = useState(false);
  const translate = useTranslate();

  useEffect(() => {
    const beforeunload = (e: BeforeUnloadEvent) => {
      e.preventDefault();
      const confirmationMessage = '';
      e.returnValue = confirmationMessage;
      return confirmationMessage;
    };

    if (messageInfo?.notificationOptions?.undoable) {
      window.addEventListener('beforeunload', beforeunload);
    }

    if (notifications.length && !messageInfo) {
      // Set a new snack when we don't have an active one
      setMessageInfo(takeNotification());
      setOpen(true);
    } else if (notifications.length && messageInfo && open) {
      // Close an active snack when a new one is added
      setOpen(false);
    }
  }, [notifications, takeNotification, open, messageInfo]);

  useEffect(() => {
    if (messageInfo) {
      const {
        message,
        type: typeFromMessage,
        notificationOptions: {
          autoHideDuration: autoHideDurationFromMessage,
          messageArgs,
          multiLine: multilineFromMessage,
          undoable,
          ...options
        },
      } = messageInfo;

      toast(translate(message, messageArgs), {
        type: typeFromMessage,
        position: toast.POSITION.BOTTOM_CENTER,
      });
    }
  }, [messageInfo]);

  return null;
};
