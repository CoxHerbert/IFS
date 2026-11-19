export type ProfileAvatarKey = 'cat' | 'bear' | 'bunny' | 'fox' | 'panda' | 'frog' | 'koala' | 'tiger' | 'lion' | 'penguin' | 'whale' | 'unicorn' | 'hamster' | 'deer' | 'owl' | 'dog' | 'chick'

export const defaultProfileAvatar = '/static/ledger-logo.svg'

export const profileAvatars: { key: ProfileAvatarKey; label: string; src: string }[] = [
  { key: 'cat', label: '小猫', src: '/static/avatars/cat.svg' },
  { key: 'bear', label: '小熊', src: '/static/avatars/bear.svg' },
  { key: 'bunny', label: '兔兔', src: '/static/avatars/bunny.svg' },
  { key: 'fox', label: '狐狸', src: '/static/avatars/fox.svg' },
  { key: 'panda', label: '熊猫', src: '/static/avatars/panda.svg' },
  { key: 'frog', label: '青蛙', src: '/static/avatars/frog.svg' },
  { key: 'koala', label: '考拉', src: '/static/avatars/koala.svg' },
  { key: 'tiger', label: '小虎', src: '/static/avatars/tiger.svg' },
  { key: 'lion', label: '狮子', src: '/static/avatars/lion.svg' },
  { key: 'penguin', label: '企鹅', src: '/static/avatars/penguin.svg' },
  { key: 'whale', label: '鲸鱼', src: '/static/avatars/whale.svg' },
  { key: 'unicorn', label: '独角兽', src: '/static/avatars/unicorn.svg' },
  { key: 'hamster', label: '仓鼠', src: '/static/avatars/hamster.svg' },
  { key: 'deer', label: '小鹿', src: '/static/avatars/deer.svg' },
  { key: 'owl', label: '猫头鹰', src: '/static/avatars/owl.svg' },
  { key: 'dog', label: '小狗', src: '/static/avatars/dog.svg' },
  { key: 'chick', label: '小鸡', src: '/static/avatars/chick.svg' },
]

export function getProfileAvatarSource(key?: string) {
  return key ? profileAvatars.find(avatar => avatar.key === key)?.src || defaultProfileAvatar : defaultProfileAvatar
}
