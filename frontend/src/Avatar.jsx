export default function Avatar({ name, size = "md" }) {
  return (
    <img
      src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${name}`}
      alt="avatar"
      className={`avatar avatar-${size}`}
    />
  );
}
