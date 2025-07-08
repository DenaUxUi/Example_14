#!/bin/bash

set -e  # Остановить скрипт при ошибке

echo "🚀 Шаг 1: Инициализация и запуск Terraform"
cd terraform
terraform init
terraform apply -auto-approve

echo "📥 Шаг 2: Получение IP из Terraform output"
IP=$(terraform output -raw instance_public_ip)
cd ../ansible
echo "$IP" > ip.txt

echo "📋 Шаг 3: Генерация Ansible inventory.ini"
cat > inventory.ini <<EOF
[ec2]
$IP ansible_user=ubuntu ansible_ssh_private_key_file=~/.ssh/Dena.pem
EOF

echo "🔧 Шаг 4: Настройка EC2 через Ansible (установка Docker и т.п.)"
ansible-playbook -i inventory.ini playbook.yml

# (необязательно, если у тебя уже автоматический запуск контейнера)
if [ -f deploy.yml ]; then
  echo "🐳 Шаг 5: Запуск контейнера через Ansible"
  ansible-playbook -i inventory.ini deploy.yml
fi

echo "✅ Готово! Приложение развернуто на IP: $IP"

