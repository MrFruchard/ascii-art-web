# Générateur d'Art ASCII Web

Ce projet est une application web écrite en Go qui permet de générer de l'art ASCII à partir d'un texte donné, en utilisant différentes bannières (`shadow`, `standard`, `thinkertoy`). Vous pouvez télécharger l'art ASCII généré sous forme de fichier texte.

## Fonctionnalités

- Génération d'art ASCII à partir de chaînes de texte.
- Utilisation de différents styles de bannières : `shadow`, `standard`, et `thinkertoy`.
- Vérification des caractères non imprimables.
- Téléchargement de l'art ASCII en fichier texte.
- Gestion des erreurs avec des pages HTML personnalisées.

## Installation et Lancement

### Pré-requis

- [Go](https://golang.org/dl/) doit être installé sur votre machine.
- Un fichier contenant les modèles ASCII pour chaque bannière (`shadow.txt`, `standard.txt`, `thinkertoy.txt`).

### Étapes d'installation

1. Clonez le dépôt sur votre machine locale :

    ```bash
    git clone https://github.com/votre-repo/generateur-ascii-art.git
    cd generateur-ascii-art
    ```

2. Placez les fichiers de bannières (`shadow.txt`, `standard.txt`, `thinkertoy.txt`) dans le répertoire racine du projet.

3. Compilez le programme :

    ```bash
    go build
    ```

4. Lancez le serveur :

    ```bash
    go run main.go
    ```

5. Ouvrez votre navigateur et accédez à l'adresse suivante :

    ```bash
    http://localhost:8044
    ```

## Utilisation

L'application dispose de deux principales fonctionnalités :

### 1. Générer et afficher l'art ASCII

- Accédez à la page principale.
- Saisissez votre texte et choisissez une bannière (ex. : `standard`, `shadow`, `thinkertoy`).
- Cliquez sur le bouton "Générer" pour afficher l'art ASCII directement dans le navigateur.

### 2. Télécharger l'art ASCII

- Après avoir généré l'art ASCII, vous pouvez le télécharger sous forme de fichier texte.
- Allez sur la page `/download` pour télécharger le fichier.


## Commandes possibles

- **Lancer le serveur :** `go run main.go`
- **Accéder à l'interface web :** `http://localhost:8044`

### Points d'entrée dans l'application

- `/` : Page principale pour générer et afficher l'art ASCII.
- `/download` : Télécharge l'art ASCII généré sous forme de fichier texte.

## Gestion des erreurs

Le programme gère les erreurs courantes :
- **Texte manquant** : Une erreur est affichée si aucun texte n'est fourni.
- **Caractères non imprimables** : Si des caractères non imprimables sont détectés, un message d'erreur est affiché.
- **Bannière manquante** : Une erreur est levée si une bannière non valide est fournie.

## Auteur

Ce projet a été développé en Go dans le cadre d'un générateur d'art ASCII avec une interface web.
