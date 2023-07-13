import logo from '../../assets/logo.jpg';

export const Logo = ({ width }: { width: string }) => {
  return (
    <div style={{ width: width }}>
      <img src={logo} style={{ display: 'block' }} width="100%" />
    </div>
  );
};
